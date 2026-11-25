package dnfbridge

import (
	"context"
	"fmt"

	"longheng.io/server/internal/modules/dnf/dnfenum"
)

// partyEntryActorBatchEnabled 已转正为默认批量版：进本前将全队成员的
// mode0（全部成员行）+ mode1（全部成员行）合并为两条包一次下发。
func partyEntryActorBatchEnabled() bool {
	return true
}

// sendCurrentPartyEntryActorBindingBatch 把整队成员的 actor 绑定合成两条包：
// 先一条 mode0（所有人），再一条 mode1（所有人）。
//
// 与支线（轻装上阵）的 party_entry_actor_current.go 同形；差别只在本服没有
// owner 替换模型，两个 owner 都取副本场景上下文（currentSceneObjectContext）。
func (s *Service) sendCurrentPartyEntryActorBindingBatch(receiver *gameSession, source string) error {
	if receiver == nil || receiver.selectedCharacterID == 0 {
		return fmt.Errorf("party entry actor batch requires selected character")
	}
	state := s.runtimePartyStateForSession(receiver)
	members := runtimePartyMembers(state)
	if len(members) < 2 || len(members) > 4 {
		return fmt.Errorf("party entry actor batch roster size=%d", len(members))
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
			actorCtx, actor, id, name, character, found, mode0Owner,
		)
		if err != nil {
			return err
		}
		summary := s.currentAccountAdventureGroupSummaryForPacket(actorCtx, actor, character, found)
		row1 := s.buildCurrentActorBindingMode1BodyForSelectedWithEquipmentInContext(
			actorCtx, actor, character, found, id, uint32(summary.ManageLevel), currentSceneObjectContext, true,
		)
		if len(row0) < 3 || len(row1) < 3 || row0[0] != 0 || row1[0] != 1 {
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
	s.logGameEvent(receiver, "game-party-entry-all-actors-bound-batch",
		"source", source, "char_id", receiver.selectedCharacterID,
		"party_id", state.PartyID, "actor_ids", actorIDs,
		"owner_channel", currentSceneObjectContext,
		"local_mode0_owner", localMode0Owner,
		"mode0_bytes", len(mode0), "mode1_bytes", len(mode1),
		"sequence", "mode0_all_then_mode1_all")
	return nil
}
