package dnfbridge

import (
	"os"
	"strings"
	"time"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
	"longheng.io/server/internal/modules/dnf/party"
)

// 重入同一 party id 的客户端 P2P 会话不会重置：现场两人打完回城后再拉退图那位，
// 服务端名册已经 member_count=3，但老成员的两条新腿一个数据报都不发，重入方一直
// "连接中"；手动踢掉一个队员（成员变化）后才恢复。
//
// 这里做的不是"修客户端"，而是服务端在检测到"该角色最近从这支队离开、现在又加回来"
// 时，把整支队退役并重建为一个**新的 party id**：客户端拿到的是全新的队伍身份和
// 新的中继端点，等价于现场那次"踢人后重新组"的状态变化。
const partyRejoinRebuildEnv = "DNFBRIDGE_PARTY_REJOIN_REBUILD"

const recentPartyDepartureTTL = 10 * time.Minute

type recentPartyDeparture struct {
	partyID uint16
	at      time.Time
}

// partyRejoinRebuildEnabled 默认开启；显式 0/off/false/legacy 可回退。
func partyRejoinRebuildEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(partyRejoinRebuildEnv))) {
	case "0", "off", "false", "legacy":
		return false
	}
	return true
}

func (s *Service) notePartyDeparture(characterID, partyID uint16) {
	if s == nil || characterID == 0 || partyID == 0 {
		return
	}
	s.recentPartyMu.Lock()
	if s.recentPartyByCharacter == nil {
		s.recentPartyByCharacter = make(map[uint16]recentPartyDeparture)
	}
	s.recentPartyByCharacter[characterID] = recentPartyDeparture{partyID: partyID, at: time.Now()}
	s.recentPartyMu.Unlock()
}

func (s *Service) departedFromPartyRecently(characterID, partyID uint16) bool {
	if s == nil || characterID == 0 || partyID == 0 {
		return false
	}
	s.recentPartyMu.Lock()
	departure, ok := s.recentPartyByCharacter[characterID]
	if ok && time.Since(departure.at) > recentPartyDepartureTTL {
		delete(s.recentPartyByCharacter, characterID)
		ok = false
	}
	s.recentPartyMu.Unlock()
	return ok && departure.partyID == partyID
}

func (s *Service) clearPartyDeparture(characterID uint16) {
	if s == nil || characterID == 0 {
		return
	}
	s.recentPartyMu.Lock()
	delete(s.recentPartyByCharacter, characterID)
	s.recentPartyMu.Unlock()
}

// rebuildRuntimePartyOnRejoin 把当前队伍原样搬到新的 party id：成员、队长、
// 各成员的会话代次不变；只有 party id（以及随之重建的中继腿）变化。
func (s *Service) rebuildRuntimePartyOnRejoin(partyID int, joinerID uint16) (alignedcmd.PartyState, bool) {
	if s == nil || partyID <= 0 || partyID > int(^uint16(0)) {
		return alignedcmd.PartyState{}, false
	}
	manager := s.runtimePartyManagerForService()
	if manager == nil {
		return alignedcmd.PartyState{}, false
	}
	snapshot, found := manager.SnapshotByID(uint16(partyID))
	if !found || len(snapshot.Members) < 2 {
		return alignedcmd.PartyState{}, false
	}
	var leader party.RuntimePartyMember
	for _, member := range snapshot.Members {
		if member.UserID == snapshot.Leader {
			leader = member
			break
		}
	}
	if leader.UserID == 0 {
		return alignedcmd.PartyState{}, false
	}
	if !manager.RetireSnapshot(snapshot) {
		return alignedcmd.PartyState{}, false
	}
	created := manager.EnsureLeader(leader, snapshot.Settings)
	if !created.OK {
		s.logWarn("party rejoin rebuild failed", "old_party_id", partyID, "stage", "ensure_leader", "reason", created.Reason)
		return alignedcmd.PartyState{}, false
	}
	newPartyID := created.Party.ID
	joined := 1
	for _, member := range snapshot.Members {
		if member.UserID == leader.UserID {
			continue
		}
		result := manager.Join(newPartyID, member)
		if !result.OK {
			s.logWarn("party rejoin rebuild failed", "old_party_id", partyID, "new_party_id", newPartyID,
				"stage", "join", "char_id", member.UserID, "reason", result.Reason)
			continue
		}
		joined++
	}
	rebuilt, found := manager.SnapshotByID(newPartyID)
	if !found || len(rebuilt.Members) != joined {
		return alignedcmd.PartyState{}, false
	}
	// 旧 party id 的中继腿必须立刻释放，新的 Sync 会按新 id 重新分配端口。
	s.closeRuntimePartyUDPRelay(partyID)
	return rebuilt.StateFor(joinerID), true
}
