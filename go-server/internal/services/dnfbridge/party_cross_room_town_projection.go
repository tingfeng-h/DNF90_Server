package dnfbridge

import (
	"longheng.io/server/internal/modules/dnf/alignedcmd"
)

// Graduated production policies:
//   - Always skip the redundant post-op11 town actor refresh in finishManagedPartyJoin
//     so P2P handshakes are never reset mid-flight.
//   - Always project same-channel Seria-room party peers with mode0+mode1 (no op23)
//     so HUD/HP/P2P work without rendering the peer inside the private Seria room.
func partyPostSnapshotTownRefreshSkipEnabled() bool {
	return true
}

func partySeriaActorNoOp23Enabled() bool {
	return true
}

// townPrivateRoomPeerProjectionPolicy decides how a same-physical-Seria-map
// party peer should be projected:
//   - blockAll=true: legacy behavior, skip both mode0+mode1 and op23.
//   - skipOp23=true: send mode0+mode1 so the peer enters the active scene table
//     (0x056A3758) for HP/party/P2P, but omit op23 so the peer is not rendered
//     inside the receiver's private Seria room.
func (s *Service) townPrivateRoomPeerProjectionPolicy(receiver *gameSession, actor onlinePlayerInfo) (blockAll bool, skipOp23 bool) {
	if !s.townPrivateRoomPeerProjectionBlocked(receiver, actor) {
		return false, false
	}
	if partySeriaActorNoOp23Enabled() {
		return false, true
	}
	return true, false
}

// refreshPostSnapshotPartyActorsAfterJoin runs at the end of
// finishManagedPartyJoin (after op11/op153/op9). When the skip switch is active,
// town receivers are NOT re-projected a second time (they were already projected
// before op11 in step 3), while active-dungeon receivers keep their idempotent
// post-snapshot projection path untouched.
func (s *Service) refreshPostSnapshotPartyActorsAfterJoin(state alignedcmd.PartyState, source string) {
	if s == nil || s.onlinePlayers == nil || state.PartyID <= 0 || len(runtimePartyMembers(state)) < 2 {
		return
	}
	if !partyPostSnapshotTownRefreshSkipEnabled() {
		s.refreshTownPartyActorsAfterJoin(state, source)
		return
	}
	dungeonReceivers := make([]*gameSession, 0, len(runtimePartyMembers(state)))
	for _, member := range runtimePartyMembers(state) {
		receiver, ok := s.onlineGameSession(member.UserID)
		if !ok || receiver == nil {
			continue
		}
		if s.sessionHasActiveDungeonRuntime(receiver) {
			dungeonReceivers = append(dungeonReceivers, receiver)
		}
	}
	if len(dungeonReceivers) == 0 {
		return
	}
	for _, receiver := range dungeonReceivers {
		if !s.recordCurrentDungeonPartyRunKey(receiver) {
			s.logGameEvent(receiver, "game-dungeon-party-run-key-after-join-failed",
				"source", source,
				"party_id", state.PartyID,
				"char_id", receiver.selectedCharacterID,
				"reason", "active_dungeon_run_or_scene_unavailable")
		}
	}
	for _, receiver := range dungeonReceivers {
		if err := s.projectDungeonPartyPeerActors(receiver, source+"_dungeon_peer_actor"); err != nil {
			s.logGameEvent(receiver, "game-dungeon-party-actor-after-join-failed",
				"source", source,
				"party_id", state.PartyID,
				"char_id", receiver.selectedCharacterID,
				"error", err)
		}
	}
}
