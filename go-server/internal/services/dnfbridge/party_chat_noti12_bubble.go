package dnfbridge

import (
	"context"
)

// Graduated production policy: party chat always emits NOTI 12 (0x000C) with
// the sender UID so both blue party chat and overhead bubbles render.
func partyChatNoti12BubbleEnabled() bool {
	return true
}

// buildCurrentPartyChatNoti12Body builds a NOTI 12 (0x000C) party chat packet
// with the real sender actor UID so DNF.exe 0x02035370 -> 0x02A04000 ->
// 0x029F54E0 displays BOTH the blue party chat line AND the overhead chat
// bubble above the speaker's character (NOTI 65 omits sender UID, causing
// FindActiveByUID(0) to return NULL and skip overhead bubbles).
func buildCurrentPartyChatNoti12Body(serverGroup byte, senderID uint16, message []byte) []byte {
	var w packetWriter
	w.writeByte(currentPartyChatType)
	w.writeUint16(senderID)
	w.writeByte(serverGroup)
	w.writeRawDstr(message)
	w.writeRawDstr(nil)
	return w.bytes()
}

func currentPartyChatNoticePayload(serverGroup byte, senderID uint16, name, message []byte) (uint16, []byte) {
	if partyChatNoti12BubbleEnabled() && senderID != 0 {
		return currentAreaChatNotifyOpcode, buildCurrentPartyChatNoti12Body(serverGroup, senderID, message)
	}
	return currentOfficialWhisperNotifyOpcode, buildCurrentPartyChatBody(serverGroup, name, message)
}

func (s *Service) sessionInActiveDungeon(session *gameSession) bool {
	if session == nil {
		return false
	}
	session.dungeon.mu.Lock()
	inDungeon := session.dungeon.runtime != nil
	session.dungeon.mu.Unlock()
	return inDungeon
}

// handleCurrentSoloDungeonPartyChat echoes type-2 chat sent by a solo player
// inside a dungeon (where the client's normal chat box automatically emits
// MessageType=2 instead of 3) via NOTI 12 so both chat log and overhead bubble
// appear instead of silently dropping the message.
func (s *Service) handleCurrentSoloDungeonPartyChat(
	ctx context.Context,
	session *gameSession,
	op uint16,
	identity boundGameSessionCharacter,
	serverGroup byte,
	name []byte,
	message []byte,
) error {
	validate := func() error {
		if !s.boundGameSessionCharacterCurrent(identity) || !s.sessionInActiveDungeon(session) {
			return errInvalidGameSessionEvent
		}
		return nil
	}
	if err := s.sendCurrentGuildChatFrameChecked(ctx, session, op, []byte{1}, 1, validate); err != nil {
		return err
	}
	noticeOpcode, body := currentPartyChatNoticePayload(serverGroup, identity.character, name, message)
	if err := s.sendCurrentGuildChatFrameChecked(ctx, session, noticeOpcode, body, 0, validate); err != nil {
		return err
	}
	s.logGameEvent(session, "game-solo-dungeon-chat-echoed",
		"source_char_id", identity.character,
		"message_len", len(message),
		"msg_id", noticeOpcode)
	return nil
}
