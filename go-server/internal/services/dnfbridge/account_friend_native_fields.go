package dnfbridge

import (
	"context"
	"os"
	"strconv"
	"strings"

	"longheng.io/server/internal/modules/dnf/accountfriend"
	repo "longheng.io/server/internal/modules/dnf/repository"
)

const (
	// DNFBRIDGE_ACCOUNT_FRIEND_NATIVE_FIELDS fixes the three NOTI1371 (op1371)
	// field mismatches verified against V212 DNF.exe (0x0201FE50 / 0x010BE12C /
	// 0x010BE8DC):
	//   1. +0x04 (NeutralState) is directory_server (channelAdvertiseID, e.g. 1),
	//      passed with +0x52 (channel) to 0x01B378F0 to render the channel name
	//      (writing 0 caused "未知").
	//   2. +0x0C (1st DSTR, cap 30) is CharacterName (rendered as "Lv.90 %s",
	//      used for right-click lookup 0x0296D810 and cross-channel op804 follow).
	//   3. +0x54 (2nd DSTR, cap 17) is AdventureName (rendered as "(%s)", used
	//      as the primary key for op1438 delete, op1440 note, and NOTI1373).
	// Default: on; set off/0/false/no/legacy to restore the legacy wire layout.
	accountFriendNativeFieldsEnv = "DNFBRIDGE_ACCOUNT_FRIEND_NATIVE_FIELDS"
)

func accountFriendNativeFieldsEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(accountFriendNativeFieldsEnv))) {
	case "off", "0", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

func (s *Service) accountFriendDirectoryServerForChannel(channel uint16) byte {
	if s == nil || channel == 0 || !accountFriendNativeFieldsEnabled() {
		return 0
	}
	dir := s.channelAdvertiseID()
	if dir <= 0 || dir >= 250 {
		return 0
	}
	return byte(dir)
}

func accountFriendNativeRowDeleteKey(row currentAccountFriendNativeRow) string {
	if accountFriendNativeFieldsEnabled() {
		return row.AdventureName
	}
	return row.CharacterName
}

func accountFriendNativeRowWireNames(r currentAccountFriendNativeRow) (first string, firstCap int, second string, secondCap int) {
	if accountFriendNativeFieldsEnabled() {
		return r.CharacterName, 30, r.AdventureName, 17
	}
	return r.AdventureName, 30, r.CharacterName, 17
}

func matchAccountFriendMutationTarget(entry accountfriend.Entry, targetName string) bool {
	return accountFriendNativeFieldsEnabled() &&
		entry.Relation.State == accountfriend.Friend &&
		entry.Account.RepresentAccountName == targetName
}

// primeSameChannelAccountFriends silently pushes a single mode8 cache frame
// (via the non-reentrant-lock-free sendCurrentCrossSceneProfileRefreshFrame)
// for each online same-channel cross-room friend in the account friend list,
// so right-clicking them in the Adventure Friend panel immediately opens the
// same-channel menu with avatar & adventure name.
func (s *Service) primeSameChannelAccountFriends(ctx context.Context, viewer boundGameSessionCharacter, rows []currentAccountFriendNativeRow) {
	if s == nil || !accountFriendNativeFieldsEnabled() ||
		s.options.currentClientVersion != "official-24.9" ||
		viewer.session == nil || !s.boundGameSessionCharacterCurrent(viewer) {
		return
	}
	viewerChannel := currentTownPresenceChannelID(viewer.session)
	if viewerChannel <= 0 || viewerChannel >= 255 {
		return
	}
	repos, ok := s.repositoryGroupFor(ctx)
	if !ok || repos.Character == nil {
		return
	}
	for _, row := range rows {
		if int(row.Channel) != viewerChannel || row.CharacterName == "" {
			continue
		}
		idStr, exists, err := repos.Character.FindIDByName(ctx, row.CharacterName)
		if err != nil || !exists {
			continue
		}
		parsedID, err := strconv.ParseUint(idStr, 10, 16)
		if err != nil || parsedID == 0 || uint16(parsedID) == viewer.character {
			continue
		}
		peer, online := s.onlineGameSessionCharacterSnapshot(uint16(parsedID))
		if !online || peer.session == nil || peer.session == viewer.session {
			continue
		}
		if currentTownPresenceChannelID(peer.session) != viewerChannel {
			continue
		}
		if s.onlinePlayers != nil && s.onlinePlayers.PeerInSameArea(viewer.character, peer.character) {
			continue
		}
		_, _ = s.primeSingleSameChannelGuildMode8(ctx, viewer, peer, viewerChannel)
	}
}

var _ = repo.CharacterRecord{}
