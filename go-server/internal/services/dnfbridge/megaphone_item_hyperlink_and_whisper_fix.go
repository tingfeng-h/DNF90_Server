package dnfbridge

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
	"strconv"
	"strings"

	"longheng.io/server/internal/modules/dnf/inventory"
	r "longheng.io/server/internal/modules/dnf/repository"
)

const (
	mainWhisperZeroTargetIDEnv  = "DNFBRIDGE_MAIN_WHISPER_ZERO_TARGET_ID"
	currentMegaphoneHyperlinkOp = uint16(372) // V212 0x02A014D0 (NOTI 372)
)

// Graduated production policy: megaphone item hyperlinks (NOTI 372) are always enabled.
func megaphoneItemHyperlinkEnabled() bool {
	return true
}

func mainWhisperZeroTargetIDEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(mainWhisperZeroTargetIDEnv))) {
	case "off", "0", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

// allowMainWhisperZeroTargetID returns true for V212 main-input whispers
// (0x02A04D66 hardcodes TargetID=0 with type=1 and trailing server_group).
func allowMainWhisperZeroTargetID(messageType byte, hasTrailing bool, targetID uint16) bool {
	return mainWhisperZeroTargetIDEnabled() && messageType == 1 && hasTrailing && targetID == 0
}

// allowMainWhisperTrailingGroup accepts the exact character server_group_id as
// well as V212 Guild/directory menu values (channelAdvertiseID, 0xFF unknown,
// or the target's owner channel byte) while preserving strict rejection in unit
// tests where channelAdvertiseID == 0 and a mismatched group is injected.
func (s *Service) allowMainWhisperTrailingGroup(mainWhisper bool, targetSession *gameSession, trailing byte, targetGroup int64) bool {
	if int64(trailing) == targetGroup {
		return true
	}
	if !mainWhisper || !mainWhisperZeroTargetIDEnabled() {
		return false
	}
	if trailing == 0xff {
		return true
	}
	if adv := s.channelAdvertiseID(); adv > 0 && adv < 255 && trailing == byte(adv) {
		return true
	}
	if ch, ok := currentSessionOwnerChannelField(targetSession); ok && ch > 0 && trailing == ch {
		return true
	}
	return false
}

// decodeCurrentMegaphoneItemHyperlinkRequest decodes op424
// (CmdPacketItemHyperlinkMessage) for channel (11) and server (13) megaphones:
//
//	u8 mode(11|13), u16 hint(0), u32 slot, DSTR message, u8 count(1..5), count*148 item bytes.
func decodeCurrentMegaphoneItemHyperlinkRequest(body []byte) (currentItemHyperlinkRequest, error) {
	if len(body) < 12 {
		return currentItemHyperlinkRequest{}, errCurrentItemHyperlinkBodyInvalid
	}
	req := currentItemHyperlinkRequest{
		MessageType: body[0],
		Hint:        binary.LittleEndian.Uint16(body[1:3]),
		Reserved:    binary.LittleEndian.Uint32(body[3:7]),
	}
	if (req.MessageType != 11 && req.MessageType != 13) || req.Hint != 0 || req.Reserved > math.MaxInt16 {
		return currentItemHyperlinkRequest{}, errCurrentItemHyperlinkBodyInvalid
	}
	msgLen := binary.LittleEndian.Uint32(body[7:11])
	if msgLen == 0 || msgLen > 1016 {
		return currentItemHyperlinkRequest{}, errCurrentItemHyperlinkBodyInvalid
	}
	msgEnd := 11 + int(msgLen)
	if msgEnd >= len(body) {
		return currentItemHyperlinkRequest{}, errCurrentItemHyperlinkBodyInvalid
	}
	req.Message = append([]byte(nil), body[11:msgEnd]...)
	if !validCurrentMegaphoneText(req.Message) {
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

// buildCurrentMegaphoneItemHyperlinkBody matches V212 NOTI 372 (0x02A014D0):
// mode(u8), server_group(u8), decoration(u8), owner_channel(u8), sender_uid(u16),
// name(DSTR), message(DSTR), count(u8) + count*148 items, group_name(DSTR).
func buildCurrentMegaphoneItemHyperlinkBody(mode, group, channel byte, id uint16, name, message []byte, items [][]byte) []byte {
	var w packetWriter
	w.writeBytes([]byte{mode, group, 0, channel})
	w.writeUint16(id)
	w.writeRawDstr(name)
	w.writeRawDstr(message)
	w.writeByte(byte(len(items)))
	for _, item := range items {
		w.writeBytes(item)
	}
	w.writeRawDstr(nil)
	return w.bytes()
}

func (s *Service) handleCurrentMegaphoneItemHyperlink(p *gameSession, req gameplayRequest) error {
	if p == nil {
		return errInvalidGameSessionEvent
	}
	ctx, cancel := guildChatContext(p)
	defer cancel()
	reject := func(err error) error {
		s.logGameEvent(p, "game-megaphone-hyperlink-rejected", "reason", err)
		return s.sendCurrentGuildChatFrame(ctx, p, req.Opcode, []byte{0, 1}, 1)
	}
	q, err := decodeCurrentMegaphoneItemHyperlinkRequest(req.Body)
	if err != nil || s.options.currentClientVersion != "official-24.9" || req.Sequence == nil || bytes.IndexByte(q.Message, 0) >= 0 {
		return reject(inventory.ErrMegaphone)
	}
	high, known := p.guild.purchaseSequence.currentHigh()
	if !known || *req.Sequence != high {
		return reject(inventory.ErrMegaphoneReplay)
	}
	identity, bound := s.boundGameSessionCharacterSnapshot(p)
	if !bound || p.loginRevoked.Load() || p.townSceneReadyCharacterID != identity.character {
		return reject(errInvalidGameSessionEvent)
	}
	if s.sessionInActiveDungeon(p) {
		return reject(inventory.ErrMegaphone)
	}
	account := s.accountIDForSession(p)
	reservation, err := s.reserveCurrentCharacterAssets(ctx, identity, account)
	if err != nil {
		return reject(err)
	}
	defer s.releaseCurrentCharacterAssets(reservation)
	repos, available := s.repositoryGroup()
	if !available || repos.Character == nil {
		return reject(inventory.ErrOwnerUnavailable)
	}
	c, found, err := repos.Character.Load(ctx, strconv.Itoa(int(identity.character)))
	if err != nil || !found || c.AccountID != account {
		return reject(inventory.ErrCharacterRequired)
	}
	name, validName := currentPartyChatName(c.Name)
	channel, validChannel := currentSessionOwnerChannelField(p)
	group := c.Stats["server_group_id"]
	if !validName || !validChannel || group < 0 || group >= 250 {
		return reject(inventory.ErrMegaphone)
	}
	catalog, err := s.currentPVFItemCatalog()
	if err != nil {
		return reject(err)
	}
	d, err := currentMegaphoneDefinition(catalog, q.MessageType)
	if err != nil {
		return reject(err)
	}
	nonce, err := p.guild.purchaseSequence.identity(0)
	if err != nil {
		return reject(err)
	}
	owner, err := inventory.NewOwner(repos)
	if err != nil {
		return reject(err)
	}
	result, err := owner.ConsumeMegaphone(ctx, inventory.MegaphoneCommand{
		AccountID: account, CharacterID: identity.character, Mode: q.MessageType, Slot: int16(q.Reserved),
		Session: nonce, Sequence: *req.Sequence, Now: s.gameplayNow().UTC(), Definition: d,
		Revalidate: func() bool { return !p.loginRevoked.Load() && s.boundGameSessionCharacterCurrent(identity) },
	})
	if err != nil {
		return reject(err)
	}
	body := buildCurrentMegaphoneItemHyperlinkBody(q.MessageType, byte(group), channel, identity.character, name, q.Message, q.Items)
	peers := s.currentMegaphonePeers(identity, q.MessageType)
	sourceChannel := p.channel.ID
	validate := func() error {
		if p.loginRevoked.Load() || !s.boundGameSessionCharacterCurrent(identity) {
			return errInvalidGameSessionEvent
		}
		return nil
	}
	sourceErr := s.sendCurrentGuildChatFrameChecked(ctx, p, req.Opcode, []byte{1}, 1, validate)
	if sourceErr == nil {
		sourceErr = s.sendCurrentGuildChatFrameChecked(ctx, p, currentMegaphoneHyperlinkOp, body, 0, validate)
	}
	if sourceErr == nil {
		entry := currentItemListEntryFromStack(r.MainInventoryListType, int16(q.Reserved), result.Stack)
		if result.Remaining == 0 {
			entry = currentItemListEntry{}
			entry.patchCore(int16(q.Reserved), math.MaxUint32, 0)
		}
		sourceErr = s.sendCurrentGuildChatFrameChecked(ctx, p, 14, buildCurrentItemUpdateBody(r.MainInventoryListType, []currentItemListEntry{entry}), 0, validate)
	}
	s.releaseCurrentCharacterAssets(reservation)
	for _, peer := range peers {
		if err := s.queueCurrentMegaphoneWithOpcode(peer, q.MessageType, sourceChannel, group, currentMegaphoneHyperlinkOp, body); err != nil {
			s.logWarn("megaphone hyperlink recipient unavailable", "target_char_id", peer.character, "error", err)
		}
	}
	s.logGameEvent(p, "game-megaphone-hyperlink-committed", "item_id", d.ItemID, "slot", q.Reserved, "remaining", result.Remaining, "mode", q.MessageType, "items", len(q.Items), "recipient_candidates", len(peers)+1)
	return sourceErr
}
