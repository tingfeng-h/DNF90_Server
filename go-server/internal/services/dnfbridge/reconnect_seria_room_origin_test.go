package dnfbridge

import (
	"bytes"
	"context"
	"encoding/binary"
	"strings"
	"testing"

	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
	dnfrepomemory "longheng.io/server/internal/modules/dnf/repository/memory"
)

// 现场缺陷：频道重连已经落在共享赛丽亚房间（38/1），但房间被当成艾尔文防线：
//   - 没有 op1375 ⇒ 客户端世界地图上下文仍是 38；
//   - 没有 townSeriaRoomEntrySnapshot ⇒ op1425 走房间所在城镇的固定外门口
//     38/0 (622,196)，而不是切换频道前所在城镇的房门。
//
// 根因是 DNFBRIDGE_RECONNECT_SERIA_ORIGIN=false 同时关掉了"落房间"和"恢复房间来源"。
// 本文件钉住拆分后的语义：
//   - 落点仍由 DNFBRIDGE_RECONNECT_SERIA_ORIGIN 决定（false = 落房间，保持本服现状）；
//   - 房间来源由 DNFBRIDGE_RECONNECT_SERIA_ENTRY_TOWN 独立恢复（默认 on）。

func reconnectSeriaRoomFixture(
	t *testing.T,
	entryTown, entryArea, entryX, entryY int64,
) (*Service, *gameSession, dnfrepo.Group) {
	t.Helper()
	service, session, repositories := newChannelReconnectTest(t)
	character, found, err := repositories.Character.Load(context.Background(), "29")
	if err != nil || !found {
		t.Fatalf("load reconnect character found=%t err=%v", found, err)
	}
	character.Stats["town_id"] = entryTown
	character.Stats["area_id"] = entryArea
	character.Stats["pos_x"] = entryX
	character.Stats["pos_y"] = entryY
	character.Stats["direction"] = 5
	character.Stats["area_state"] = 3
	// 刻意删除持久化 prev-village，覆盖"老角色第一次切频道"的回退路径：
	// 来源只能从 commit 之前的角色记录（切换前位置）推出来。
	for _, key := range []string{
		currentTownPrevVillageTownStat,
		currentTownPrevVillageAreaStat,
		currentTownPrevVillagePosXStat,
		currentTownPrevVillagePosYStat,
		currentTownPrevVillageDirectionStat,
	} {
		delete(character.Stats, key)
	}
	// commitCurrentReconnectLocation writes through CharacterSettlement; the
	// default fake store and the memory group's settlement transaction own
	// different character rows. Put the character in one aggregate and use its
	// character/settlement repositories, as the PVP reconnect fixtures do.
	aggregate := dnfrepomemory.NewMemoryGroup()
	if err := aggregate.Character.Save(context.Background(), character); err != nil {
		t.Fatal(err)
	}
	repositories.Character = aggregate.Character
	repositories.CharacterAssets = aggregate.CharacterAssets
	repositories.CharacterSettlement = aggregate.CharacterSettlement
	// newTownMoveTest's provider closure captured its own copy of the group;
	// re-bind it so the settlement transaction and the loader see this row.
	service.repositoryProvider = func() (dnfrepo.Group, bool) { return repositories, true }
	service.seriaRoomDoorMapReader = fakePVFMapReader(t, map[string]string{
		"WestCoast.map":     westCoastMapPVFText,
		"new_Elvengard.map": elvengardMapPVFText,
	})

	session.selectedCharacterID = 29
	session.accountID = "dnf:1"
	session.channelReconnect = true
	service.bindGameSessionCharacter(session, 29)
	session.residentChannel = session.channel
	session.connectionTownActorOwnerChannel = byte(session.channel.ID)
	session.townActorOwnerChannel = byte(session.channel.ID)
	return service, session, repositories
}

func findWorldTownContextPacket(t *testing.T, packetsBytes []byte) ([]byte, bool) {
	t.Helper()
	for _, packet := range splitAllCurrentGameServerUpperPackets(t, packetsBytes) {
		if packet.Header.MsgID == currentTownWorldTownContextNotificationMsgID {
			return packet.Body, true
		}
	}
	return nil, false
}

// (a) 修复档（默认 ENTRY_TOWN=on）：频道重连落 38/1，但房间来源 = 切换前城镇 40/0。
// 出房间（op1425）回到西海岸主街的赛丽亚房门 (523,142)，而不是 38/0 (622,196)。
func TestChannelReconnectSeriaRoomKeepsThePreSwitchTown(t *testing.T) {
	t.Setenv(reconnectSeriaOriginEnv, "false")
	t.Setenv(reconnectSeriaEntryTownEnv, "")
	service, session, repositories := reconnectSeriaRoomFixture(t, 40, 0, 1500, 500)
	readLog := newRecordedPacketLog(t, service)

	if err := service.sendCurrentChannelReconnectTownEntry(session); err != nil {
		t.Fatalf("channel reconnect must land in the room: %v", err)
	}
	// ★ 本次修复的绑定点：落点是赛丽亚房间 ⇒ 房间来源 = 切换前城镇 40/0。
	if entry := session.townSeriaRoomEntrySnapshot; !entry.PositionValid ||
		entry.CharacterID != 29 || entry.TownID != 40 || entry.AreaID != 0 ||
		entry.PositionX != 1500 || entry.PositionY != 500 {
		t.Fatalf("reconnect must bind the room origin to the pre-switch town 40/0 (1500,500), got %+v", entry)
	}
	if origin := session.townPrevVillageSnapshot; !origin.PositionValid ||
		origin.TownID != 40 || origin.AreaID != 0 || origin.PositionX != 1500 || origin.PositionY != 500 {
		t.Fatalf("prev-village snapshot must be the pre-switch town 40/0 (1500,500), got %+v", origin)
	}

	// 落点仍是共享房间 38/1，且来源已经随同一个 settlement 事务落库。
	stored := loadStoredCharacter(t, repositories)
	if stored.Stats["town_id"] != 38 || stored.Stats["area_id"] != 1 {
		t.Fatalf("durable landing must stay the shared room 38/1: %+v", stored.Stats)
	}
	if stored.Stats[currentTownPrevVillageTownStat] != 40 ||
		stored.Stats[currentTownPrevVillageAreaStat] != 0 ||
		stored.Stats[currentTownPrevVillagePosXStat] != 1500 ||
		stored.Stats[currentTownPrevVillagePosYStat] != 500 {
		t.Fatalf("room origin was not persisted with the reconnect landing: %+v", stored.Stats)
	}

	// op1375 的世界地图归属城镇必须是 40，而不是房间所在的 38。
	worldBody, found := findWorldTownContextPacket(t, session.conn.(*bufferConn).write.Bytes())
	if !found {
		t.Fatalf("missing op1375 world-town context; client would treat the room as Elvenguard")
	}
	if len(worldBody) != 4 || binary.LittleEndian.Uint32(worldBody) != 40 {
		t.Fatalf("op1375 world town body=%x want 40", worldBody)
	}

	// 第一次走下方门：来源城镇 40 的房门 (523,142)，不是 38/0 (622,196)。
	body := lastSeriaRoomExitBody(t, service, session)
	if want := wantSeriaRoomExitBody(t, 40, 0, 523, 142, 5, 3); !bytes.Equal(body, want) {
		t.Fatalf("first room exit body=%x want the pre-switch town's door 40/0 (523,142) body=%x", body, want)
	}
	returned := loadStoredCharacter(t, repositories)
	if returned.Stats["town_id"] != 40 || returned.Stats["area_id"] != 0 ||
		returned.Stats["pos_x"] != 523 || returned.Stats["pos_y"] != 142 {
		t.Fatalf("first room exit must return to the pre-switch town 40/0 (523,142): %+v", returned.Stats)
	}

	text := readLog()
	for _, want := range []string{
		"game-channel-reconnect-town-entry-send",
		"seria_origin_bound=true",
		"seria_origin_town_id=40",
		"seria_origin_area_id=0",
		"seria_entry_town_bound=true",
		"landing_source=seria_room_entry_town_door",
		"entry_used=true",
		"entry_town_id=40",
		"door_map_path=WestCoast.map",
		"door_trigger_x=473",
		"door_trigger_y=127",
		"door_trigger_width=100",
		"door_trigger_height=30",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("fixed reconnect-then-exit log missing %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "landing_source=seria_room_outer_door") {
		t.Fatalf("fixed behaviour still answered the room town's outer door:\n%s", text)
	}
}

// (b) 回退档 ENTRY_TOWN=off：落点仍是房间，但来源一位都不绑 ⇒ 逐字节回到线上错行为
// （op1425 = 38/0 (622,196)，且没有 op1375 世界城镇通知）。
func TestChannelReconnectSeriaRoomOriginSwitchOffKeepsTheOuterDoor(t *testing.T) {
	t.Setenv(reconnectSeriaOriginEnv, "false")
	t.Setenv(reconnectSeriaEntryTownEnv, "off")
	service, session, repositories := reconnectSeriaRoomFixture(t, 40, 0, 1500, 500)
	readLog := newRecordedPacketLog(t, service)

	if err := service.sendCurrentChannelReconnectTownEntry(session); err != nil {
		t.Fatal(err)
	}
	if entry := session.townSeriaRoomEntrySnapshot; entry.PositionValid {
		t.Fatalf("ENTRY_TOWN=off must not bind the room origin, got %+v", entry)
	}
	if origin := session.townPrevVillageSnapshot; origin.PositionValid {
		t.Fatalf("ENTRY_TOWN=off must not bind prev-village, got %+v", origin)
	}
	if _, found := findWorldTownContextPacket(t, session.conn.(*bufferConn).write.Bytes()); found {
		t.Fatal("ENTRY_TOWN=off must not send the op1375 world-town context")
	}

	body := lastSeriaRoomExitBody(t, service, session)
	if want := wantSeriaRoomExitBody(t, 38, 0, 622, 196, 5, 3); !bytes.Equal(body, want) {
		t.Fatalf("rollback exit body=%x want today's outer door 38/0 (622,196) body=%x", body, want)
	}
	returned := loadStoredCharacter(t, repositories)
	if returned.Stats["town_id"] != 38 || returned.Stats["area_id"] != 0 ||
		returned.Stats["pos_x"] != 622 || returned.Stats["pos_y"] != 196 {
		t.Fatalf("rollback must persist the outer door 38/0 (622,196): %+v", returned.Stats)
	}

	text := readLog()
	for _, want := range []string{
		"game-reconnect-seria-origin-skipped",
		"landing_source=seria_room_outer_door",
		"entry_used=false",
		"seria_entry_resolve_reason=no_seria_room_entry_origin_for_this_character",
		"door_town_id=38",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("rollback log missing %q:\n%s", want, text)
		}
	}
}

// (c) 持久化 prev-village 优先：老角色已经有 40/0 的持久化来源时，即使当前记录
// 已经是 38/1 房间，也恢复 40/0（不依赖本次 pre-commit 位置）。
func TestChannelReconnectSeriaRoomPrefersPersistedPrevVillage(t *testing.T) {
	t.Setenv(reconnectSeriaOriginEnv, "false")
	t.Setenv(reconnectSeriaEntryTownEnv, "")
	service, session, repositories := reconnectSeriaRoomFixture(t, 38, 1, 450, 234)
	character := loadStoredCharacter(t, repositories)
	currentTownPrevVillageApplyOriginStats(character.Stats, currentTownPositionSnapshot{
		CharacterID: 29, TownID: 40, AreaID: 0, PositionX: 1500, PositionY: 500, MovementCode: 5, PositionValid: true,
	})
	if err := repositories.Character.Save(context.Background(), character); err != nil {
		t.Fatal(err)
	}

	if err := service.sendCurrentChannelReconnectTownEntry(session); err != nil {
		t.Fatal(err)
	}
	if entry := session.townSeriaRoomEntrySnapshot; !entry.PositionValid ||
		entry.TownID != 40 || entry.AreaID != 0 || entry.PositionX != 1500 || entry.PositionY != 500 {
		t.Fatalf("persisted prev-village 40/0 (1500,500) must win over the 38/1 session position, got %+v", entry)
	}
}

// (d) 切频道前刚走到的城镇优先于数据库里尚未回写的旧 prev-village。
// 现场：角色从 14 工业基地走到 20 苍穹贵族号，立刻切频道；DB 仍是 14，
// 但 session 已经在 20。频道重连必须绑定 20，而不是把房间送回 14。
func TestChannelReconnectSeriaRoomPrefersLiveSessionTownOverStaleDB(t *testing.T) {
	t.Setenv(reconnectSeriaOriginEnv, "false")
	t.Setenv(reconnectSeriaEntryTownEnv, "")
	service, session, repositories := reconnectSeriaRoomFixture(t, 40, 0, 1500, 500)

	// DB 的持久化来源仍是旧城镇 40/0；session 已经站在 39/0。
	character := loadStoredCharacter(t, repositories)
	currentTownPrevVillageApplyOriginStats(character.Stats, currentTownPositionSnapshot{
		CharacterID: 29, TownID: 40, AreaID: 0, PositionX: 1500, PositionY: 500, MovementCode: 5, PositionValid: true,
	})
	if err := repositories.Character.Save(context.Background(), character); err != nil {
		t.Fatal(err)
	}
	session.townPositionSnapshot = currentTownPositionSnapshot{
		CharacterID: 29, TownID: 39, AreaID: 0, PositionX: 596, PositionY: 144, MovementCode: 5, PositionValid: true,
	}
	session.townSceneReadyCharacterID = 29
	service.seriaRoomDoorMapReader = fakePVFMapReader(t, map[string]string{
		"WestCoast.map":     westCoastMapPVFText,
		"new_Elvengard.map": elvengardMapPVFText,
		"new_HendonMyre.map": "[town movable area with animation]\n[movable area]\n[area]\n" +
			"546 129 100 30\n[move town]\n38 1\n[/movable area]\n[/town movable area with animation]\n",
	})

	if err := service.sendCurrentChannelReconnectTownEntry(session); err != nil {
		t.Fatal(err)
	}
	if entry := session.townSeriaRoomEntrySnapshot; !entry.PositionValid ||
		entry.TownID != 39 || entry.AreaID != 0 || entry.PositionX != 596 || entry.PositionY != 144 {
		t.Fatalf("live session town 39/0 (596,144) must win over stale DB origin 40/0, got %+v", entry)
	}
	stored := loadStoredCharacter(t, repositories)
	if stored.Stats[currentTownPrevVillageTownStat] != 39 ||
		stored.Stats[currentTownPrevVillageAreaStat] != 0 {
		t.Fatalf("the live session town must be persisted with the landing: %+v", stored.Stats)
	}
	body := lastSeriaRoomExitBody(t, service, session)
	if want := wantSeriaRoomExitBody(t, 39, 0, 596, 144, 5, 3); !bytes.Equal(body, want) {
		t.Fatalf("exit body=%x want the live session town's door 39/0 (596,144) body=%x", body, want)
	}
}

// (e) 新连接（session 快照为空）时，角色记录里**当前** town/area 优先于
// 尚未更新的 prev_village_*。现场等价于：DB 已在 20/1，prev_village 仍是 14/0。
func TestChannelReconnectSeriaRoomPrefersCharacterCurrentTownOverStalePrevVillage(t *testing.T) {
	t.Setenv(reconnectSeriaOriginEnv, "false")
	t.Setenv(reconnectSeriaEntryTownEnv, "")
	service, session, repositories := reconnectSeriaRoomFixture(t, 40, 0, 1500, 500)
	session.townPositionSnapshot = currentTownPositionSnapshot{}
	session.townSceneReadyCharacterID = 0

	character := loadStoredCharacter(t, repositories)
	currentTownPrevVillageApplyOriginStats(character.Stats, currentTownPositionSnapshot{
		CharacterID: 29, TownID: 39, AreaID: 0, PositionX: 596, PositionY: 144, MovementCode: 5, PositionValid: true,
	})
	if err := repositories.Character.Save(context.Background(), character); err != nil {
		t.Fatal(err)
	}

	if err := service.sendCurrentChannelReconnectTownEntry(session); err != nil {
		t.Fatal(err)
	}
	if entry := session.townSeriaRoomEntrySnapshot; !entry.PositionValid ||
		entry.TownID != 40 || entry.AreaID != 0 || entry.PositionX != 1500 || entry.PositionY != 500 {
		t.Fatalf("character current town 40/0 (1500,500) must win over stale prev_village 39/0, got %+v", entry)
	}
	stored := loadStoredCharacter(t, repositories)
	if stored.Stats[currentTownPrevVillageTownStat] != 40 ||
		stored.Stats[currentTownPrevVillageAreaStat] != 0 {
		t.Fatalf("character current town must be persisted with the landing: %+v", stored.Stats)
	}
	body := lastSeriaRoomExitBody(t, service, session)
	if want := wantSeriaRoomExitBody(t, 40, 0, 523, 142, 5, 3); !bytes.Equal(body, want) {
		t.Fatalf("exit body=%x want the current town's door 40/0 (523,142) body=%x", body, want)
	}
}

// (f) 登录落点开关不能改变“切频道必进房间”：DNFBRIDGE_LOGIN_PERSISTED_TOWN=keep
// 会让共享登录解析器返回持久化街道位置（例如 40/0），频道重连必须把它改成 38/1
// 赛丽亚房间，而不是“从哪条街切就回哪条街”。
func TestChannelReconnectSeriaRoomIgnoresLoginPersistedTownKeep(t *testing.T) {
	t.Setenv(reconnectSeriaOriginEnv, "false")
	t.Setenv(reconnectSeriaEntryTownEnv, "")
	t.Setenv(currentLoginPersistedTownEnv, currentLoginPersistedTownKeepValue)
	service, session, repositories := reconnectSeriaRoomFixture(t, 40, 0, 1500, 500)

	if err := service.sendCurrentChannelReconnectTownEntry(session); err != nil {
		t.Fatal(err)
	}
	stored := loadStoredCharacter(t, repositories)
	if stored.Stats["town_id"] != 38 || stored.Stats["area_id"] != 1 {
		t.Fatalf("LOGIN_PERSISTED_TOWN=keep must not make a channel switch land in the street: %+v", stored.Stats)
	}
	if entry := session.townSeriaRoomEntrySnapshot; !entry.PositionValid ||
		entry.TownID != 40 || entry.AreaID != 0 || entry.PositionX != 1500 || entry.PositionY != 500 {
		t.Fatalf("room origin must still be the pre-switch town 40/0 (1500,500), got %+v", entry)
	}
	body := lastSeriaRoomExitBody(t, service, session)
	if want := wantSeriaRoomExitBody(t, 40, 0, 523, 142, 5, 3); !bytes.Equal(body, want) {
		t.Fatalf("exit body=%x want the pre-switch town's door 40/0 (523,142) body=%x", body, want)
	}
}

// 开关默认档：未设置/空/未知值 = 修复档；只有显式回退值才关闭。
func TestChannelReconnectSeriaEntryTownDefaultsToTheFix(t *testing.T) {
	for _, value := range []string{"", "   ", "on", "true", "1", "yes", "fixed", "garbage"} {
		t.Run("enabled_"+strings.ReplaceAll(value, " ", "_"), func(t *testing.T) {
			t.Setenv(reconnectSeriaEntryTownEnv, value)
			if !currentReconnectSeriaEntryTownEnabled() {
				t.Fatalf("%s=%q must stay on the fixed default tier", reconnectSeriaEntryTownEnv, value)
			}
		})
	}
	for _, value := range []string{"off", "OFF", " off ", "legacy", "Legacy", "0", "false", "no"} {
		t.Run("legacy_"+strings.ReplaceAll(value, " ", "_"), func(t *testing.T) {
			t.Setenv(reconnectSeriaEntryTownEnv, value)
			if currentReconnectSeriaEntryTownEnabled() {
				t.Fatalf("%s=%q must fall back to the legacy tier", reconnectSeriaEntryTownEnv, value)
			}
		})
	}
}
