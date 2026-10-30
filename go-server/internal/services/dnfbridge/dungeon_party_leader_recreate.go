package dnfbridge

import (
	"longheng.io/server/internal/modules/dnf/alignedcmd"
	"longheng.io/server/internal/modules/dnf/party"
)

// \u961f\u957f\u9000\u56fe\u3001\u4e14\u961f\u4f0d\u5269\u4f59 \u22652 \u4eba\u65f6\uff1a\u5ba2\u6237\u7aef\u5728\u526f\u672c\u5185**\u4e0d\u91cd\u5efa\u961f\u4f0d\u8868**\uff08
// op9 \u53ea\u4f1a\u66f4\u65b0\u5df2\u6709\u884c\uff0c\u800c\u4e14\u65e7\u961f\u957f\u4f1a\u88ab\u6807\u6210\u79bb\u7ebf"\u7ea2\u7535\u8111"\uff09\uff0c\u4f46"\u89e3\u6563"\u8fd9\u6761\u8def\u5b83\u662f\u8ba4\u7684\u3002
// \u6240\u4ee5\u8fd9\u91cc\u7528\u201c\u5148\u89e3\u6563\u3001\u518d\u7528\u5269\u4f59\u73a9\u5bb6\u91cd\u5efa\u4e00\u652f\u65b0\u961f\u4f0d\u201d\u6765\u6a21\u62df\u91cd\u5efa\uff1a
// \u5ba2\u6237\u7aef\u4f1a\u628a\u65b0\u961f\u4f0d\u5f53\u6210\u4e00\u6b21\u5168\u65b0\u7684\u7ec4\u961f\uff08\u65b0 party_id\u3001\u65b0\u961f\u957f\u5728 slot 0\uff09\u3002
//
// \u4efb\u4f55\u4e00\u6b65\u5931\u8d25\u90fd\u8fd4\u56de ok=false\uff0c\u7531\u8c03\u7528\u65b9\u56de\u9000\u5230\u666e\u901a\u8f6c\u8ba9/\u89e3\u6563\u3002
func (s *Service) recreateRuntimePartyOnDungeonLeaderExit(
	session *gameSession,
	instance *dungeonPartyInstance,
	leavingCharacterID uint16,
	source string,
) (int, bool) {
	if s == nil || session == nil || leavingCharacterID == 0 {
		return 0, false
	}
	state := s.runtimePartyStateForSession(session)
	members := runtimePartyMembers(state)
	if state.PartyID <= 0 || len(members) < 3 || state.UserID != leavingCharacterID {
		return 0, false
	}
	manager := s.runtimePartyManagerForService()
	if manager == nil {
		return 0, false
	}
	// \u5e78\u5b58\u8005\u6309\u69fd\u4f4d\u987a\u5e8f\uff0c\u7b2c\u4e00\u4e2a\u5c31\u662f\u65b0\u961f\u957f\u3002
	var survivors []boundGameSessionCharacter
	var survivorMembers []party.RuntimePartyMember
	for _, member := range members {
		if member.UserID == 0 || member.UserID == leavingCharacterID {
			continue
		}
		identity, online := s.onlineGameSessionCharacterSnapshot(member.UserID)
		if !online || identity.session == nil {
			return 0, false
		}
		pm, ok := s.runtimePartyMemberForIdentity(identity)
		if !ok {
			return 0, false
		}
		survivors = append(survivors, identity)
		survivorMembers = append(survivorMembers, pm)
	}
	if len(survivors) < 2 {
		return 0, false
	}
	identity, bound := s.boundGameSessionCharacterSnapshot(session)
	if !bound || !s.boundGameSessionCharacterCurrent(identity) {
		return 0, false
	}

	// 1) \u89e3\u6563\u539f\u961f\u4f0d\uff08\u5148\u975e\u961f\u957f\u3001\u6700\u540e\u961f\u957f\uff0c\u907f\u514d\u89e6\u53d1\u201c\u961f\u957f\u79bb\u961f\u2192\u91cd\u5efa\u4e00\u4ee3\u201d\uff09\u3002
	for _, member := range members {
		if member.UserID == 0 || member.UserID == leavingCharacterID {
			continue
		}
		manager.LeaveCurrent(member.UserID)
	}
	manager.LeaveCurrent(leavingCharacterID)

	// 2) \u5148\u628a\u6bcf\u4e2a\u5e78\u5b58\u8005\u5ba2\u6237\u7aef\u7684\u65e7\u961f\u4f0d\u8868\u6e05\u6389\uff08\u7a7a\u540d\u518c = \u62c6\u8868\uff09\u3002
	for _, identity := range survivors {
		if err := s.sendRuntimePartyMemberRemoved(identity.session, 0, alignedcmd.PartyState{}); err != nil {
			s.logGameEvent(identity.session, "game-party-dungeon-leader-exit-recreate-clear-failed",
				"source", source, "char_id", identity.character, "error", err)
		}
		s.clearRuntimePartyProjection(identity.character, identity.generation)
	}

	// 3) \u7528\u5e78\u5b58\u8005\u91cd\u5efa\u4e00\u652f\u65b0\u961f\u4f0d\uff1a\u69fd\u4f4d\u6700\u9760\u524d\u7684\u4eba\u5f53\u961f\u957f\u3002
	settings := runtimePartyStateSnapshot(survivors[0].session)
	created := manager.EnsureLeader(survivorMembers[0], settings)
	if !created.OK {
		s.logGameEvent(session, "game-party-dungeon-leader-exit-recreate-failed",
			"source", source, "old_party_id", state.PartyID,
			"leader_char_id", survivors[0].character, "reason", created.Reason)
		return 0, false
	}
	for _, member := range survivorMembers[1:] {
		if result := manager.Join(created.Party.ID, member); !result.OK {
			s.logGameEvent(session, "game-party-dungeon-leader-exit-recreate-join-failed",
				"source", source, "party_id", created.Party.ID,
				"char_id", member.UserID, "reason", result.Reason)
			return 0, false
		}
	}
	snapshot, found := manager.SnapshotByID(created.Party.ID)
	if !found {
		return 0, false
	}
	newState := s.publishRuntimePartySnapshot(snapshot)
	s.syncRuntimePartyUDPRelay(newState)

	// 4) \u6309\u201c\u521a\u7ec4\u961f\u201d\u7684\u5b8c\u6574\u6295\u5f71\u53d1\u7ed9\u6240\u6709\u4eba\u3002
	if err := s.sendManagedRuntimePartySnapshots(newState); err != nil {
		s.logGameEvent(session, "game-party-dungeon-leader-exit-recreate-snapshot-failed",
			"source", source, "party_id", newState.PartyID, "error", err)
	}
	s.sendDungeonPartyFrameRoster(newState, source+"_recreated_dungeon_roster")

	// 5) \u65e7\u961f\u957f\u672c\u4eba\uff1a\u57ce\u9547\u4e0a\u4e0b\u6587\u9500\u6bc1\u5e27 + \u9000\u961f ACK\u3002
	if err := s.sendDungeonExitSelfFrameClear(session, source+"_leader"); err != nil {
		s.logGameEvent(session, "game-party-dungeon-leader-exit-recreate-self-clear-failed", "source", source, "error", err)
	}
	if err := s.sendDungeonExitPartyLeaveAck(session, source+"_leader"); err != nil {
		s.logGameEvent(session, "game-party-dungeon-leader-exit-recreate-self-ack-failed", "source", source, "error", err)
	}
	s.clearRuntimePartyProjection(leavingCharacterID, identity.generation)

	s.logGameEvent(session, "game-party-dungeon-leader-exit-recreated",
		"source", source,
		"old_party_id", state.PartyID,
		"new_party_id", newState.PartyID,
		"leaving_char_id", leavingCharacterID,
		"new_leader_id", survivors[0].character,
		"member_count", len(runtimePartyMembers(newState)),
		"instance_id", instanceIDOfPartyDungeonInstance(instance),
		"reason", "client_only_updates_dungeon_rows_so_the_party_is_dissolved_and_recreated_with_the_survivors")
	return len(survivors), true
}
