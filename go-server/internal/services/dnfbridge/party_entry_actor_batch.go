package dnfbridge

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"longheng.io/server/internal/modules/dnf/dnfenum"
)

// 移植自「轻装上阵」分支 party_entry_actor_current.go 的思路：
// 进副本加载前（op28/op29 之前）把**全队** actor 一次性绑定——
// 先发一条 mode0（含全部队员行），再发一条 mode1（含全部队员行）。
//
// 该分支的现场结论：加载完成之后补发的队友投影，无法替换房间在绑定参与者时
// 已经消费掉的 actor。我们此前所有"后补发"实验失败，可能正是时机晚了一步。
//
//	DNFBRIDGE_PARTY_ENTRY_BATCH_ACTORS=off   # 默认开启；置 off 回退到逐人补发
//
// 实机验收（2026-10-29）：party_id=1 首进与 party_id=2 退队重组后均绑定全队，
// 双方模型互相可见，日志无失败。
const partyEntryBatchActorsEnv = "DNFBRIDGE_PARTY_ENTRY_BATCH_ACTORS"

func partyEntryBatchActorsEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(partyEntryBatchActorsEnv))) {
	case "0", "off", "false", "no":
		return false
	default:
		return true
	}
}

// currentLocalActorMode0OwnerContext 给出本地角色在副本场景里 mode0 应使用的 owner。
// 非副本场景沿用调用方给的默认值；副本里用会话已提交的频道，无法用字节表示时明确失败。
func currentLocalActorMode0OwnerContext(session *gameSession, existingOwner byte) (byte, error) {
	if !currentDungeonSceneActive(session) {
		return existingOwner, nil
	}
	owner := currentClientSelectedChannel(session)
	if owner == currentSceneObjectContext {
		return 0, fmt.Errorf(
			"local dungeon mode0 owner is not representable: selected_channel=%d character_id=%d",
			currentClientSelectedChannelID(session),
			session.selectedCharacterID,
		)
	}
	return owner, nil
}

// sendCurrentPartyEntryActorBinding 把整队 actor 一次性绑定到收包者。
func (s *Service) sendCurrentPartyEntryActorBinding(receiver *gameSession, source string) error {
	state := s.runtimePartyStateForSession(receiver)
	members := runtimePartyMembers(state)
	if len(members) < 2 || len(members) > 4 {
		return fmt.Errorf("party entry actor roster size=%d", len(members))
	}
	ctx, cancel := context.WithTimeout(context.Background(), createWriteTimeout)
	defer cancel()
	localMode0Owner, err := currentLocalActorMode0OwnerContext(receiver, currentSceneObjectContext)
	if err != nil {
		return err
	}
	mode0 := []byte{0, byte(len(members)), 0}
	mode1 := []byte{1, byte(len(members)), 0}
	actorIDs := make([]uint16, 0, len(members))
	// 先把所有行都构造好再发，避免队友缺失时留下"看起来完整"的半套名册。
	for _, member := range members {
		actor, online := s.onlineGameSession(member.UserID)
		if !online || actor == nil || actor.selectedCharacterID != member.UserID {
			return fmt.Errorf("party entry actor %d is unavailable", member.UserID)
		}
		actorCtx := ctx
		if actor != receiver {
			actorCtx = withDungeonPartyPeerSkillPageTarget(ctx, receiver.selectedCharacterID)
		}
		id, name, character, found := s.characterForEnter(actorCtx, actor, member.UserID)
		if !found || id != member.UserID {
			return fmt.Errorf("party entry actor %d record is unavailable", member.UserID)
		}
		mode0Owner := byte(currentSceneObjectContext)
		if actor == receiver {
			mode0Owner = localMode0Owner
		}
		row0, err := s.buildCurrentSceneObjectListBodyForSessionInContextStrict(
			actorCtx, actor, id, name, character, found, mode0Owner)
		if err != nil {
			return err
		}
		summary := s.currentAccountAdventureGroupSummaryForPacket(actorCtx, actor, character, found)
		row1 := s.buildCurrentActorBindingMode1BodyForSelectedWithEquipmentInContext(
			actorCtx, actor, character, found, id, uint32(summary.ManageLevel), currentSceneObjectContext, true)
		if len(row0) < 3 || len(row1) < 3 || row0[0] != 0 || row1[0] != 1 ||
			binary.LittleEndian.Uint16(row0[1:3]) != 1 || binary.LittleEndian.Uint16(row1[1:3]) != 1 {
			return fmt.Errorf("party entry actor %d has invalid typed rows", id)
		}
		mode0 = append(mode0, row0[3:]...)
		mode1 = append(mode1, row1[3:]...)
		actorIDs = append(actorIDs, id)
	}
	for _, body := range [][]byte{mode0, mode1} {
		if err := s.sendGameUpperRawClass(receiver, uint16(dnfenum.CmdPacketSetUDPIPPort), body, 0); err != nil {
			return err
		}
	}
	receiver.preDungeonContextPlayerStateSent = true
	s.logGameEvent(receiver, "game-party-entry-all-actors-bound",
		"source", source, "char_id", receiver.selectedCharacterID,
		"party_id", state.PartyID, "actor_ids", actorIDs, "owner_channel", currentSceneObjectContext,
		"local_mode0_owner_channel", localMode0Owner,
		"mode0_bytes", len(mode0), "mode1_bytes", len(mode1),
		"sequence", "mode0_all_then_mode1_all_before_op28_op29")
	return nil
}
