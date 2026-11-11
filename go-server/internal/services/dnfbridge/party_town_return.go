package dnfbridge

import (
	"fmt"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
)

// rebuildRuntimePartyTownCoPresenceAfterReturn restores the scene-owned town
// actors only after every same-area party member has consumed its own op132
// scene finalizer.  The typed op24 return packet intentionally contains just
// the receiving character, so its old peer actors have been detached even
// though OnlinePlayerManager still correctly records the party in that town.
// Replaying a peer before it consumes the finalizer corrupts its selector;
// waiting for all members gives both clients the native mode0/mode1/op9/op23
// reconstruction at the same point.
func (s *Service) rebuildRuntimePartyTownCoPresenceAfterReturn(session *gameSession, source string) error {
	if s == nil || session == nil || session.selectedCharacterID == 0 || s.onlinePlayers == nil {
		return nil
	}
	state := s.townCopresenceActorPartyState(session)
	members := runtimePartyMembers(state)
	// A one-member party also needs realtime actor state on the dungeon-to-town
	// transition too. Its peer loop is empty, but the existing roster rebuild
	// must still run after the same canonical party and town-scene checks.
	if state.PartyID <= 0 || len(members) == 0 {
		return nil
	}

	sessions := make([]*gameSession, 0, len(members))
	for _, member := range members {
		memberSession, online := s.onlineGameSession(member.UserID)
		if !online || memberSession == nil || !currentPartyTownReturnSceneReady(memberSession) {
			return nil
		}
		memberState := s.townCopresenceActorPartyState(memberSession)
		if memberState.PartyID != state.PartyID || memberState.UserID != state.UserID ||
			len(runtimePartyMembers(memberState)) != len(members) {
			return nil
		}
		sessions = append(sessions, memberSession)
	}

	rebuildSig, claimed := claimPartyTownReturnRebuild(state, sessions)
	if !claimed {
		return nil
	}

	for _, receiver := range sessions {
		receiverChannelID := currentTownPresenceChannelID(receiver)
		aliasRetired, err := s.retireDungeonPartyAliasRecordsIfPending(receiver, members, source)
		if err != nil {
			releasePartyTownReturnRebuild(sessions, rebuildSig)
			return err
		}
		for _, member := range members {
			if member.UserID == receiver.selectedCharacterID {
				continue
			}
			actor, found := s.onlinePlayers.PlayerForCharacter(member.UserID)
			if !found || actor.Session == nil {
				releasePartyTownReturnRebuild(sessions, rebuildSig)
				return fmt.Errorf("party return peer %d has no town actor", member.UserID)
			}
			// 频道隔离（彻底隔离口径）：跨频道队友一律不投影（op23 没有频道字段，
			// 发它等于把对方声明成本地区域的一员；同 sendTownRemotePartyActors）。
			if actor.ChannelID != receiverChannelID {
				continue
			}
			if s.townPrivateRoomPeerProjectionBlocked(receiver, actor) {
				continue
			}
			// 当 aliasRetired == true 时，上面刚发的 kind3(owner=0) 触发了客户端
			// NOTI9 的 0x2977AF0 (ClearPartyManager)，把 0x56A3768+0x48 的 slotActors
			// 槽位清空了；即便 PeerInSameArea 为 true 也必须重投一次同屏队友的
			// mode0/mode1 把城镇 actor 重新绑回槽位，否则队友血条永远读不到。
			if !aliasRetired && s.onlinePlayers.PeerInSameArea(receiver.selectedCharacterID, member.UserID) {
				continue
			}
			if err := s.sendTownRemoteActorState(receiver, actor, source+"_party_peer_actor"); err != nil {
				releasePartyTownReturnRebuild(sessions, rebuildSig)
				return err
			}
			areaBody := buildCurrentTownUserAreaNotificationBody(
				currentSceneActorObjectKey(actor.CharacterID),
				actor.TownID,
				actor.AreaID,
				actor.PositionX,
				actor.PositionY,
				actor.Direction,
				actor.AreaState,
			)
			if err := s.sendCurrentSceneFixedClass0Packet(
				receiver,
				currentTownUserAreaNotificationMsgID,
				areaBody,
				source+"_party_peer_area",
			); err != nil {
				releasePartyTownReturnRebuild(sessions, rebuildSig)
				return err
			}
		}
		receiverState := stateForRuntimePartyReceiver(state, receiver.selectedCharacterID)
		if len(members) == 1 {
			// The generic roster helper deliberately omits realtime packets for
			// single-member formation. A completed town return needs the actor
			// state restored, so this path uses the established nonempty realtime body
			// before restoring its actor frame. Keep formation behavior unchanged.
			if err := s.sendGameUpperRawClassCodec(receiver, 0x0099, runtimePartyRealtimeInfoBody(state), 0, true); err != nil {
				releasePartyTownReturnRebuild(sessions, rebuildSig)
				return err
			}
			if _, err := s.sendCurrentPartyActorFrameProjection(receiver, receiver, receiverState, source+"_party_roster"); err != nil {
				releasePartyTownReturnRebuild(sessions, rebuildSig)
				return err
			}
			if aliasRetired {
				if err := s.sendGameUpperRawClassCodec(receiver, 0x0099, runtimePartyRealtimeInfoBody(state), 0, true); err != nil {
					releasePartyTownReturnRebuild(sessions, rebuildSig)
					return err
				}
			}
			continue
		}
		if err := s.sendRuntimePartyRosterLocal(receiver, receiverState, source+"_party_roster"); err != nil {
			releasePartyTownReturnRebuild(sessions, rebuildSig)
			return err
		}
		// 当本次刚执行过 kind3(owner=0) 退役时，sendRuntimePartyRosterLocal 内部先发的
		// 那条 op153 到达客户端时 [local_actor+0x4C8] 还是 NULL（被 kind3 清空），会被
		// 0x6EC1600 判定 left_party 丢弃；等紧随其后的 op9 kind0 把绑定恢复之后，
		// 必须补发一条 op153 把血条写回。
		if aliasRetired {
			if err := s.sendRuntimePartyRealtimeInfoLocal(receiver, receiverState); err != nil {
				releasePartyTownReturnRebuild(sessions, rebuildSig)
				return err
			}
		}
	}
	s.logGameEvent(session, "game-party-town-return-copresence-rebuilt",
		"source", source,
		"party_id", state.PartyID,
		"member_count", len(members),
		"rebuild_sig", rebuildSig,
		"reason", "all_same_party_members_consumed_op132_scene_finalizer_then_mode0_mode1_op9_op23_and_realtime_roster")
	return nil
}

func currentPartyTownReturnSceneReady(session *gameSession) bool {
	if session == nil || session.selectedCharacterID == 0 {
		return false
	}
	session.dungeon.mu.Lock()
	pendingReturn := session.confirmedDungeonReturnStatePending
	hasDungeonRuntime := session.dungeon.runtime != nil
	session.dungeon.mu.Unlock()
	if pendingReturn || hasDungeonRuntime {
		return false
	}
	// op35 may finish the actor rebuild before the trailing op15 arrives.
	// Keep that op15 suppression flag armed, but a completed generation is
	// already ready for the town party/HUD restoration.
	if session.backToVillageEnterSelectPending {
		session.townPostTransition.mu.Lock()
		complete := session.townPostTransition.characterID == session.selectedCharacterID &&
			session.townPostTransition.stage == currentTownPostTransitionComplete
		session.townPostTransition.mu.Unlock()
		if !complete {
			return false
		}
	}
	session.townMu.Lock()
	ready := session.townSceneReadyCharacterID == session.selectedCharacterID &&
		session.townPositionSnapshot.CharacterID == session.selectedCharacterID &&
		session.townPositionSnapshot.PositionValid
	session.townMu.Unlock()
	return ready
}

func stateForRuntimePartyReceiver(state alignedcmd.PartyState, receiverID uint16) alignedcmd.PartyState {
	state = cloneRuntimePartyState(state)
	state.IsLeader = state.UserID == receiverID
	return state
}
