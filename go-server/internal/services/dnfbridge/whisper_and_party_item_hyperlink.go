package dnfbridge

import (
	"bytes"
	"context"
	"encoding/binary"
	"strconv"
	"strings"
)

// Graduated production policy: whisper and party item hyperlinks are always enabled.
func whisperAndPartyItemHyperlinkEnabled() bool {
	return true
}

// decodeCurrentWhisperItemHyperlinkRequest decodes V212 op424 (0x02A037A1..0x02A039E9)
// for whispers (type=1 main input, type=0x24/0x28 detached window):
//
//	u8 type, u16 targetID, u32 unknown(0), DSTR message, DSTR targetName,
//	u8 trailingGroup, u8 itemCount(1..5), itemCount*148 item bytes.
func decodeCurrentWhisperItemHyperlinkRequest(body []byte) (currentPeerChatRequest, [][]byte, error) {
	if len(body) < 16 {
		return currentPeerChatRequest{}, nil, errCurrentItemHyperlinkBodyInvalid
	}
	req := currentPeerChatRequest{
		MessageType: body[0],
		TargetID:    binary.LittleEndian.Uint16(body[1:3]),
		Unknown:     binary.LittleEndian.Uint32(body[3:7]),
	}
	if (req.MessageType != 1 && !isCurrentPrivateWhisperType(req.MessageType)) || req.Unknown != 0 {
		return currentPeerChatRequest{}, nil, errCurrentItemHyperlinkBodyInvalid
	}
	message, next, ok := readCurrentRawDstr(body, 7)
	if !ok || len(message) == 0 || len(message) > currentItemHyperlinkMessageLimit || bytes.IndexByte(message, 0) >= 0 {
		return currentPeerChatRequest{}, nil, errCurrentItemHyperlinkBodyInvalid
	}
	targetName, next, ok := readCurrentRawDstr(body, next)
	if !ok || len(targetName) == 0 {
		return currentPeerChatRequest{}, nil, errCurrentItemHyperlinkBodyInvalid
	}
	if next+2 > len(body) {
		return currentPeerChatRequest{}, nil, errCurrentItemHyperlinkBodyInvalid
	}
	req.TrailingValue = body[next]
	req.HasTrailingValue = true
	next++
	count := int(body[next])
	if count == 0 || count > currentItemHyperlinkItemLimit {
		return currentPeerChatRequest{}, nil, errCurrentItemHyperlinkBodyInvalid
	}
	itemsStart := next + 1
	itemsEnd := itemsStart + count*currentItemHyperlinkItemSize
	if itemsEnd < itemsStart || itemsEnd != len(body) {
		return currentPeerChatRequest{}, nil, errCurrentItemHyperlinkBodyInvalid
	}
	items := make([][]byte, 0, count)
	for off := itemsStart; off < itemsEnd; off += currentItemHyperlinkItemSize {
		items = append(items, append([]byte(nil), body[off:off+currentItemHyperlinkItemSize]...))
	}
	req.Message = message
	req.TargetName = targetName
	return req, items, nil
}

// buildCurrentWhisperItemHyperlinkNoti371Body matches V212 NOTI 371 (0x02A04520)
// for named whisper messages with item links (no 2-byte guild decoration).
func buildCurrentWhisperItemHyperlinkNoti371Body(messageType, group byte, sender string, message []byte, items [][]byte) []byte {
	var w packetWriter
	w.writeBytes([]byte{messageType, 0})
	w.writeRawDstr(rosterNameBytes(sender))
	w.writeByte(group)
	w.writeRawDstr(message)
	w.writeByte(byte(len(items)))
	for _, item := range items {
		w.writeBytes(item)
	}
	w.writeRawDstr(nil)
	return w.bytes()
}

func buildCurrentPartyItemHyperlinkNoti370Body(serverGroup byte, senderID uint16, message []byte, items [][]byte) []byte {
	var w packetWriter
	w.writeByte(currentPartyChatType)
	w.writeUint16(senderID)
	w.writeByte(serverGroup)
	w.writeRawDstr(message)
	w.writeByte(byte(len(items)))
	for _, item := range items {
		w.writeBytes(item)
	}
	w.writeRawDstr(nil)
	return w.bytes()
}

func decodeCurrentPartyItemHyperlinkRequest(body []byte) (currentItemHyperlinkRequest, error) {
	if len(body) < 12 {
		return currentItemHyperlinkRequest{}, errCurrentItemHyperlinkBodyInvalid
	}
	req := currentItemHyperlinkRequest{
		MessageType: body[0],
		Hint:        binary.LittleEndian.Uint16(body[1:3]),
		Reserved:    binary.LittleEndian.Uint32(body[3:7]),
	}
	if req.MessageType != currentPartyChatType || req.Reserved != 0 {
		return currentItemHyperlinkRequest{}, errCurrentItemHyperlinkBodyInvalid
	}
	msgLen := binary.LittleEndian.Uint32(body[7:11])
	if msgLen == 0 || msgLen > currentItemHyperlinkMessageLimit {
		return currentItemHyperlinkRequest{}, errCurrentItemHyperlinkBodyInvalid
	}
	msgEnd := 11 + int(msgLen)
	if msgEnd >= len(body) {
		return currentItemHyperlinkRequest{}, errCurrentItemHyperlinkBodyInvalid
	}
	req.Message = append([]byte(nil), body[11:msgEnd]...)
	if bytes.IndexByte(req.Message, 0) >= 0 {
		return currentItemHyperlinkRequest{}, errCurrentItemHyperlinkBodyInvalid
	}
	count := int(body[msgEnd])
	if count == 0 || count > currentItemHyperlinkItemLimit {
		return currentItemHyperlinkRequest{}, errCurrentItemHyperlinkBodyInvalid
	}
	itemsStart := msgEnd + 1
	itemsEnd := itemsStart + count*currentItemHyperlinkItemSize
	if itemsEnd < itemsStart || itemsEnd != len(body) {
		return currentItemHyperlinkRequest{}, errCurrentItemHyperlinkBodyInvalid
	}
	req.Items = make([][]byte, 0, count)
	for off := itemsStart; off < itemsEnd; off += currentItemHyperlinkItemSize {
		req.Items = append(req.Items, append([]byte(nil), body[off:off+currentItemHyperlinkItemSize]...))
	}
	return req, nil
}

func (s *Service) handleCurrentWhisperItemHyperlink(session *gameSession, typ uint16, body []byte) error {
	req, items, err := decodeCurrentWhisperItemHyperlinkRequest(body)
	if err != nil || s.options.currentClientVersion != "official-24.9" {
		s.logGameEvent(session, "game-whisper-hyperlink-rejected", "body_len", len(body), "error", err)
		return s.sendGameUpperFailure(session, typ, 1)
	}
	sourceID := selectedCharacterID(session)
	ctx, cancel := context.WithTimeout(context.Background(), createWriteTimeout)
	defer cancel()
	repos, available := s.repositoryGroup()
	if !available || repos.Character == nil || sourceID == 0 || strings.TrimSpace(s.accountIDForSession(session)) == "" {
		return s.sendGameUpperFailure(session, typ, 1)
	}
	mainWhisper := req.MessageType == 1
	if req.TargetID == 0 || (mainWhisper && req.TargetID == 0xffff) {
		name, decodeErr := decodeCharacterName(req.TargetName)
		if decodeErr != nil {
			return s.sendGameUpperFailure(session, typ, 1)
		}
		id, found, findErr := repos.Character.FindIDByName(ctx, name)
		value, parseErr := strconv.ParseUint(id, 10, 16)
		if findErr != nil || !found || parseErr != nil || value == 0 {
			return s.sendGameUpperFailure(session, typ, 1)
		}
		req.TargetID = uint16(value)
	}
	if req.TargetID == 0 || req.TargetID == sourceID {
		return s.sendGameUpperFailure(session, typ, 1)
	}
	targetSession, online := s.onlineGameSession(req.TargetID)
	if !online || selectedCharacterID(targetSession) != req.TargetID {
		return s.sendGameUpperFailure(session, typ, 1)
	}
	sourceIdentity, sourceBound := s.boundGameSessionCharacterSnapshot(session)
	targetIdentity, targetBound := s.boundGameSessionCharacterSnapshot(targetSession)
	source, sourceFound, sourceErr := repos.Character.Load(ctx, strconv.Itoa(int(sourceID)))
	target, targetFound, targetErr := repos.Character.Load(ctx, strconv.Itoa(int(req.TargetID)))
	targetName, nameErr := decodeCharacterName(req.TargetName)
	if !sourceBound || !targetBound || sourceErr != nil || targetErr != nil || !sourceFound || !targetFound ||
		strings.TrimSpace(source.AccountID) != strings.TrimSpace(s.accountIDForSession(session)) ||
		strings.TrimSpace(target.AccountID) != strings.TrimSpace(s.accountIDForSession(targetSession)) ||
		nameErr != nil || targetName != target.Name ||
		!s.boundGameSessionCharacterCurrent(sourceIdentity) || !s.boundGameSessionCharacterCurrent(targetIdentity) {
		return s.sendGameUpperFailure(session, typ, 1)
	}
	sourceGroup, targetGroup := source.Stats["server_group_id"], target.Stats["server_group_id"]
	if sourceGroup < 0 || sourceGroup >= 250 || targetGroup != sourceGroup ||
		!s.allowMainWhisperTrailingGroup(mainWhisper, targetSession, req.TrailingValue, targetGroup) {
		return s.sendGameUpperFailure(session, typ, 1)
	}
	if err := s.sendGameUpperSuccess(session, typ, nil); err != nil {
		return err
	}
	group := currentAdventureCharacterServerGroup(source)
	outType := req.MessageType
	if isCurrentPrivateWhisperType(outType) {
		outType = 0x24
	}
	notice := buildCurrentWhisperItemHyperlinkNoti371Body(outType, group, source.Name, req.Message, items)
	if err := s.sendGameUpperRawClassCodec(targetSession, currentGuildItemHyperlinkNotifyOpcode, notice, 0, true); err != nil {
		return err
	}
	if !mainWhisper {
		if err := s.sendGameUpperRawClassCodec(session, currentGuildItemHyperlinkNotifyOpcode, notice, 0, true); err != nil {
			return err
		}
	}
	s.logGameEvent(session, "game-whisper-hyperlink-forwarded",
		"source_char_id", sourceID, "target_char_id", req.TargetID, "message_type", req.MessageType, "items", len(items))
	return nil
}

func (s *Service) handleCurrentPartyItemHyperlink(session *gameSession, typ uint16, body []byte) error {
	ctx, cancel := guildChatContext(session)
	defer cancel()
	reject := func() error { return s.sendCurrentGuildChatFrame(ctx, session, typ, []byte{0, 1}, 1) }
	req, err := decodeCurrentPartyItemHyperlinkRequest(body)
	if err != nil {
		return reject()
	}
	identity, bound := s.boundGameSessionCharacterSnapshot(session)
	if !bound {
		return reject()
	}
	account := s.accountIDForSession(session)
	reservation, err := s.reserveCurrentCharacterAssets(ctx, identity, account)
	if err != nil {
		return reject()
	}
	defer s.releaseCurrentCharacterAssets(reservation)
	repos, available := s.repositoryGroup()
	if !available || repos.Character == nil {
		return reject()
	}
	actor, found, err := repos.Character.Load(ctx, strconv.Itoa(int(identity.character)))
	if err != nil || !found || account == "" || actor.AccountID != account {
		return reject()
	}
	group, validGroup := currentNamedChatSenderServerGroup(session, actor)
	if !validGroup {
		return reject()
	}
	notice := buildCurrentPartyItemHyperlinkNoti370Body(group, identity.character, req.Message, req.Items)
	manager := s.runtimePartyManagerForService()
	state, inParty := manager.SnapshotByUser(identity.character, identity.generation)
	if !inParty || len(state.Members) < 2 {
		if s.sessionInActiveDungeon(session) {
			validate := func() error {
				if !s.boundGameSessionCharacterCurrent(identity) || !s.sessionInActiveDungeon(session) {
					return errInvalidGameSessionEvent
				}
				return nil
			}
			if err := s.sendCurrentGuildChatFrameChecked(ctx, session, typ, []byte{1}, 1, validate); err != nil {
				return err
			}
			return s.sendCurrentGuildChatFrameChecked(ctx, session, currentItemHyperlinkNotifyOpcode, notice, 0, validate)
		}
		return reject()
	}
	validate := func() error {
		current, belongs := manager.SnapshotByUser(identity.character, identity.generation)
		if !belongs || current.ID != state.ID || len(current.Members) < 2 || !s.boundGameSessionCharacterCurrent(identity) {
			return errInvalidGameSessionEvent
		}
		return nil
	}
	if err := s.sendCurrentGuildChatFrameChecked(ctx, session, typ, []byte{1}, 1, validate); err != nil {
		return err
	}
	if err := s.sendCurrentGuildChatFrameChecked(ctx, session, currentItemHyperlinkNotifyOpcode, notice, 0, validate); err != nil {
		return err
	}
	s.releaseCurrentCharacterAssets(reservation)
	for _, member := range state.Members {
		if member.UserID == identity.character {
			continue
		}
		peer, online := s.onlineGameSessionCharacterSnapshot(member.UserID)
		if !online || peer.session.channel.ID != session.channel.ID {
			continue
		}
		current, belongs := manager.SnapshotByUser(peer.character, peer.generation)
		if !belongs || current.ID != state.ID || len(current.Members) < 2 {
			continue
		}
		_ = s.queueCurrentPartyChatWithOpcode(peer, state.ID, currentItemHyperlinkNotifyOpcode, notice)
	}
	return nil
}
