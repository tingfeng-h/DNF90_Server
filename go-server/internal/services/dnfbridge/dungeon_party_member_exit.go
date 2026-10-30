package dnfbridge

import (
	"errors"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
	"longheng.io/server/internal/modules/dnf/dnfenum"
	"longheng.io/server/internal/modules/dnf/party"
)

// A member's explicit mid-run giveup leaves the party after the town scene
// is confirmed. Completion and automatic deadline returns never call this.
// The remaining members retain their dungeon runtime and shared instance.
func (s *Service) leaveRuntimePartyOnMemberGiveup(session *gameSession, instance *dungeonPartyInstance) error {
	identity, bound := s.boundGameSessionCharacterSnapshot(session)
	if !bound || !s.boundGameSessionCharacterCurrent(identity) || instance == nil {
		return nil
	}
	manager := s.runtimePartyManagerForService()
	before, found := manager.SnapshotByUser(identity.character, identity.generation)
	drifted := !found
	if drifted {
		// Only the current authoritative socket may repair recorded generation
		// drift. An old dungeon return must never remove a new party membership.
		if !s.boundGameSessionCharacterCurrent(identity) {
			return nil
		}
		before, found = manager.SnapshotByUser(identity.character, 0)
	}
	if !found || before.ID != instance.key.PartyID || before.Leader == identity.character {
		return nil
	}
	peers := dungeonPartyInstanceMemberSessions(instance, identity.character)
	var next alignedcmd.PartyState
	var removed bool
	if s.options.currentClientVersion == "official-24.9" {
		// The leader is still inside this exact instance. Do not pass through
		// the generic below-two-members publisher: it erases the surviving
		// lobby's projection and recycles its still-active dungeon registry.
		var result party.RuntimePartyResult
		if drifted {
			result = manager.LeaveCurrent(identity.character)
		} else {
			result = manager.Leave(identity.character, identity.generation)
		}
		removed = result.OK
		if removed {
			s.logRuntimePartyMembershipRemoved("member_giveup_confirmed_town", identity.character, identity.generation, before, result)
			s.clearRuntimePartyProjection(identity.character, identity.generation)
			next = s.publishRuntimePartySnapshot(result.Party)
		}
	} else if drifted {
		next, _, removed = s.repairManagedRuntimePartyLeave(identity)
	} else {
		next, _, removed = s.leaveManagedRuntimeParty(identity, "member_giveup_confirmed_town")
	}
	if !removed {
		return nil
	}
	var errs error
	if s.options.currentClientVersion == "official-24.9" {
		errs = errors.Join(errs, s.clearDisbandedDungeonPartyRecords(session, before.StateFor(before.Leader).Members, before.Leader, false, "member_giveup_retired_receiver"))
	}
	if s.options.currentClientVersion != "official-24.9" {
		for _, peer := range peers {
			body := buildCurrentSceneOp9ActorRemovalBodyInContext(currentSceneActorObjectKey(identity.character), currentSceneObjectContext, s.currentSceneActorRoute())
			errs = errors.Join(errs, s.sendGameUpperRawClass(peer.Session, uint16(dnfenum.CmdPacketRecoverStamina), body, 0))
		}
	}
	if len(runtimePartyMembers(next)) < 2 {
		s.closeRuntimePartyUDPRelay(int(before.ID))
	} else {
		s.syncRuntimePartyUDPRelay(next)
	}
	// 踢人（walkout）路径会给包括被踢者在内的所有成员发一条
	// "某槽位已离队"的 op7 通知，客户端认这条才会把那一行从队伍界面拆掉；
	// 退图路径以前只发名册销毁帧。这里按踢人路径补上。
	if beforeMembers := runtimePartyMembers(before.StateFor(before.Leader)); len(beforeMembers) > 0 {
		for _, member := range beforeMembers {
			if member.UserID != identity.character {
				continue
			}
			if err := s.broadcastRuntimePartyWalkoutNotice(before.StateFor(before.Leader), member.Slot, 0); err != nil {
				s.logGameEvent(session, "game-party-dungeon-exit-walkout-notice-failed",
					"source", "member_giveup_confirmed_town", "char_id", identity.character, "slot", member.Slot, "error", err)
			} else {
				s.logGameEvent(session, "game-party-dungeon-exit-walkout-notice",
					"source", "member_giveup_confirmed_town", "char_id", identity.character, "slot", member.Slot,
					"members", len(beforeMembers))
			}
			break
		}
	}
	// 退图退队也要发同一条退队 ACK，否则客户端不会拆掉自己的队伍界面。
	errs = errors.Join(errs, s.sendDungeonExitPartyLeaveAck(session, "member_giveup_confirmed_town"))
	errs = errors.Join(errs, s.sendDungeonExitSelfFrameClear(session, "member_giveup_confirmed_town"))
	errs = errors.Join(errs, s.sendRuntimePartyMemberRemoved(session, identity.character, alignedcmd.PartyState{}))
	if s.options.currentClientVersion == "official-24.9" {
		errs = errors.Join(errs, s.sendMemberGiveupSurvivorRosterAfterDeparture(next, instance, identity.character))
	} else {
		errs = errors.Join(errs, s.broadcastRuntimePartyMemberRemoved(before.StateFor(before.Leader), next, identity.character, session))
	}
	s.logGameEvent(session, "game-party-member-giveup-left", "char_id", identity.character,
		"party_id", before.ID, "leader_char_id", before.Leader, "remaining_members", len(next.Members))
	return errs
}

// Rebuild only the remaining records in their current scene context. A town
// owner record cannot replace a survivor's old owner-0 dungeon roster.
func (s *Service) sendMemberGiveupSurvivorRoster(state alignedcmd.PartyState, instance *dungeonPartyInstance) error {
	return s.sendMemberGiveupSurvivorRosterAfterDeparture(state, instance, 0)
}

func (s *Service) sendMemberGiveupSurvivorRosterAfterDeparture(state alignedcmd.PartyState, instance *dungeonPartyInstance, departingID uint16) error {
	var errs error
	members := runtimePartyMembers(state)
	// The return caller holds its own dungeon lock. Never acquire a peer's
	// session lock: simultaneous returns would acquire these in opposite order.
	peers := dungeonPartyInstanceMemberSessions(instance, 0)
	for _, member := range members {
		identity, ok := s.onlineGameSessionCharacterSnapshot(member.UserID)
		if !ok {
			continue
		}
		receiver := identity.session
		active := false
		for _, peer := range peers {
			if peer.CharacterID == identity.character && peer.Session == receiver {
				active = true
				break
			}
		}
		if !active {
			errs = errors.Join(errs, s.sendRuntimePartyMemberRemoved(receiver, 0, state))
			continue
		}
		if departingID != 0 && len(members) == 1 {
			// The server keeps this player's active instance and fixed slot, but
			// a one-player dungeon has no remote party HUD/network membership.
			// Retire both old records, self last, before NOTI6 erases the peer
			// identity. Requires the client local-actor-preserving kind3 helper:
			// deleting the local party record must not destroy its living model.
			retired := []alignedcmd.PartyMemberState{{UserID: departingID}, member}
			if err := s.clearDungeonSurvivorTownPartyRecords(receiver, retired, "member_giveup_single_survivor"); err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			if err := s.clearDisbandedDungeonPartyRecords(receiver, retired, 0, false, "member_giveup_single_survivor"); err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			if err := s.sendDungeonDepartedActorRemoval(receiver, departingID, "member_giveup_single_survivor"); err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			// 队伍只剩 1 人：服务端也必须把这条队伍解散掉。只发空 op153 不清 op9，
			// 客户端会留下一条没有归属的成员行（城镇/副本里都能看到那条残留血条）；
			// party_id 继续存在还会在城镇同屏和再次进团本时被投影回来，
			// 超时空二次进图就会卡在选图/加载。先解散管理器里的队伍，再发
			// 空名册（kind3 销毁）+ 空 op153 + 清 actor 投影。
			if manager := s.runtimePartyManagerForService(); manager != nil {
				if result := manager.LeaveCurrent(identity.character); result.OK {
					s.logGameEvent(receiver, "game-party-single-survivor-disbanded",
						"char_id", identity.character,
						"departing_char_id", departingID,
						"party_id", state.PartyID,
						"disbanded", result.Disbanded,
						"reason", "a_one_member_party_must_not_survive_the_departure")
				}
			}
			s.clearRuntimePartyProjection(identity.character, identity.generation)
			_err := s.sendDungeonExitPartyLeaveAck(receiver, "member_giveup_single_survivor")
			errs = errors.Join(errs, _err)
			errs = errors.Join(errs, s.sendRuntimePartyMemberRemoved(receiver, 0, alignedcmd.PartyState{}))
			continue
		}
		if departingID != 0 {
			// First apply the new membership to this receiver's local record,
			// while the old actor slots still identify every network peer.
			sent, err := s.sendCurrentPartyActorFrameProjectionInContext(receiver, receiver, state, "member_giveup_network_roster", nil, currentPartyFrameOwnerSelection{Mode: currentPartyFrameOwnerDungeonContext})
			if err != nil || !sent {
				errs = errors.Join(errs, err)
				continue
			}
			body := buildCurrentSceneOp9ActorRemovalBodyInContext(currentSceneActorObjectKey(departingID), currentSceneObjectContext, s.currentSceneActorRoute())
			if err := s.sendGameUpperRawClass(receiver, uint16(dnfenum.CmdPacketRecoverStamina), body, 0); err != nil {
				errs = errors.Join(errs, err)
				continue
			}
		}
		var rosterErr error
		selfSent := false
		for _, actorRow := range runtimePartyMembers(state) {
			actor, found := s.onlineGameSession(actorRow.UserID)
			if !found {
				continue
			}
			sent, err := s.sendCurrentPartyActorFrameProjectionInContext(receiver, actor, state, "member_giveup_survivor_dungeon_roster", nil, currentPartyFrameOwnerSelection{Mode: currentPartyFrameOwnerDungeonContext})
			selfSent = selfSent || (sent && actor == receiver)
			rosterErr = errors.Join(rosterErr, err)
		}
		errs = errors.Join(errs, rosterErr)
		// NOTI9 diffs the old manager slots to retire missing network peers.
		// NOTI6 clears those slots before erasing the actor. Sending it first
		// loses the only identity that the roster diff can use for removal.
		// Keep each receiver's actor until its complete roster has been sent.
		if rosterErr == nil && selfSent && departingID != 0 {
			errs = errors.Join(errs, s.sendDungeonDepartedActorRemoval(receiver, departingID, "member_giveup_confirmed_town"))
		}
		// 三/四人局队伍要继续存在，不能整队销毁；但退队者本人那一行
		// 在**两套上下文**（城镇 owner=10、副本 owner=0）都要按人清掉，
		// 否则图里的队伍框会一直留着他的血条。
		if departingID != 0 {
			retired := []alignedcmd.PartyMemberState{{UserID: departingID}}
			if err := s.clearDungeonSurvivorTownPartyRecords(receiver, retired, "member_giveup_network_roster_retired"); err != nil {
				errs = errors.Join(errs, err)
			}
			if err := s.clearDisbandedDungeonPartyRecords(receiver, retired, 0, false, "member_giveup_network_roster_retired"); err != nil {
				errs = errors.Join(errs, err)
			}
		}
		errs = errors.Join(errs, s.sendGameUpperRawClassCodec(receiver, 0x0099, runtimePartyRealtimeInfoBody(state), 0, true))
	}
	return errs
}
