package dnfbridge

import (
	"context"
	"errors"
	"os"
	"strings"

	dnfguild "longheng.io/server/internal/modules/dnf/guild"
	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
)

// reconnectSeriaOriginEnv controls the channel-reconnect landing mode.
//
//	unset / on / true / 1 / yes: keep the old persisted-position landing.
//	off / false / no / 0:       land in the shared Seria room, using the same
//	                            resolver as character-list login.
//
// The name is historical: the value that used to disable the room-origin
// restore is now the switch that requests the room landing. Origin restoration
// is controlled separately by reconnectSeriaEntryTownEnv.
const reconnectSeriaOriginEnv = "DNFBRIDGE_RECONNECT_SERIA_ORIGIN"

// reconnectSeriaEntryTownEnv controls whether a room landing restores the
// shared Seria room's owning town (the native world-map context sent by op1375
// and the bottom-door origin consumed by op1425).
//
//	unset / on / true / 1 / yes: fixed behavior. Bind both snapshots to the
//	                             last non-Seria town, so leaving the room
//	                             returns to that town instead of 38/0.
//	off / legacy / false / 0 / no: legacy behavior. Leave both snapshots empty;
//	                               the client treats the room as Elvenguard and
//	                               exits to 38/0 (622,196).
//
// This is deliberately separate from DNFBRIDGE_RECONNECT_SERIA_ORIGIN: the
// production server currently sets that legacy switch to false to request the
// room landing, and doing so must not also suppress the origin fix.
const reconnectSeriaEntryTownEnv = "DNFBRIDGE_RECONNECT_SERIA_ENTRY_TOWN"

func currentReconnectSeriaEntryTownEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(reconnectSeriaEntryTownEnv))) {
	case "off", "legacy", "0", "false", "no":
		return false
	default:
		return true
	}
}

// currentReconnectSessionTownOriginLocked returns the town/area/position the
// live session is standing in right now. A channel switch loads the character
// record from the database before a just-finished cross-town walk has been
// persisted, so the durable prev_village_* row can still point at the town the
// player left. The session snapshot (op35/op36) is the authoritative "where I
// am switching from" in that window.
//
// The snapshot is deliberately not validated against character.Stats: that row
// is exactly the stale value this helper exists to override. It is validated
// against the runtime PVF and the same persistable-origin rules instead.
func (s *Service) currentReconnectSessionTownOriginLocked(session *gameSession, characterID uint16) (currentTownPositionSnapshot, bool) {
	if s == nil || session == nil || characterID == 0 {
		return currentTownPositionSnapshot{}, false
	}
	snapshot := session.townPositionSnapshot
	if snapshot.CharacterID != characterID || !snapshot.PositionValid ||
		session.townSceneReadyCharacterID != characterID {
		return currentTownPositionSnapshot{}, false
	}
	if snapshot.TownID == 0 || currentRegisteredChannelTown(int64(snapshot.TownID)) {
		return currentTownPositionSnapshot{}, false
	}
	if _, persistable := s.currentTownPrevVillageOriginIsPersistable(snapshot); !persistable {
		return currentTownPositionSnapshot{}, false
	}
	return snapshot, true
}

// 调用方持有townMu。在重连首次1375/23/24之前，沿用普通登录已校验的
// 持久化房门来源；只恢复房间的归属城镇与出口上下文，不改变房内落点或玩家资产。
//
// character must be the record **before** commitCurrentReconnectLocation
// overwrites town_id with the room landing. The candidate order is deliberate:
//
//  1. live session snapshot (same connection / transport reconnect);
//  2. the character record's **current** town/area (the walk may have updated
//     town_id but not prev_village_*);
//  3. durable prev_village_* (the last non-Seria room origin).
//
// A channel switch usually opens a new game connection, so (1) is often empty;
// (2) is then the freshest server-side evidence. The old order used
// prev_village_* first and could bind the town the player left (e.g. 14
// 工业基地) while the character record was already in 20 苍穹贵族号.
func (s *Service) restoreCurrentReconnectSeriaOriginLocked(ctx context.Context, session *gameSession, characterID uint16, character dnfrepo.CharacterRecord, mapPath string) {
	session.townPrevVillageSnapshot = currentTownPositionSnapshot{}
	session.townSeriaRoomEntrySnapshot = currentTownPositionSnapshot{}
	if !currentReconnectSeriaEntryTownEnabled() {
		s.logGameEvent(session, "game-reconnect-seria-origin-skipped",
			"char_id", characterID, "env", reconnectSeriaEntryTownEnv)
		return
	}
	if !townMapPathLooksLikeSeria(mapPath) || isPVPChannel(session.channel) {
		return
	}
	if _, raid := currentRaidLoginTownID(session); raid {
		return
	}
	sessionOrigin, _ := s.currentReconnectSessionTownOriginLocked(session, characterID)
	characterOrigin, _ := currentTownPrevVillageOriginLocked(session, characterID, character)
	candidates := []struct {
		origin currentTownPositionSnapshot
		source string
	}{
		{sessionOrigin, "channel_reconnect_session_town_position"},
		{characterOrigin, "channel_reconnect_character_current_town"},
	}
	// The persisted row is the fallback, not the primary source.
	if persisted, found := currentTownPrevVillageOriginFromPersistedStats(characterID, character); found {
		candidates = append(candidates, struct {
			origin currentTownPositionSnapshot
			source string
		}{persisted, "channel_reconnect_persisted_prev_village"})
	}
	for _, candidate := range candidates {
		origin := candidate.origin
		if !origin.PositionValid || currentRegisteredChannelTown(int64(origin.TownID)) {
			continue
		}
		originMap, persistable := s.currentTownPrevVillageOriginIsPersistable(origin)
		if !persistable {
			continue
		}
		if townMapPathLooksLikeGuildHall(originMap) {
			if _, err := s.currentTownAreaInstance(ctx, session, character, originMap); err != nil {
				if !errors.Is(err, dnfguild.ErrNotMember) {
					s.logGameEvent(session, "重连房门来源校验失败", "char_id", characterID, "error", err)
				}
				continue
			}
		}
		session.townPrevVillageSnapshot = origin
		session.townSeriaRoomEntrySnapshot = origin
		s.logTownPrevVillageOriginBound(session, origin, candidate.source)
		return
	}
}
