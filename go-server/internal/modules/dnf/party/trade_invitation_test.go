package party

import (
	"testing"
	"time"

	"longheng.io/server/internal/modules/dnf/alignedcmd"
)

func TestTradeInvitationCancelEnvelope(t *testing.T) {
	body := []byte{5, 0, 1, 0, 0, 0, 0, 133, 0}
	request, err := DecodeResponsePeerRequest(body)
	if err != nil || request.TargetID != 5 || request.Mode != 1 || request.RejectReason != 133 {
		t.Fatalf("cancel=%+v/%v", request, err)
	}
	for _, bad := range [][]byte{body[:8], append(append([]byte{}, body...), 0), {5, 0, 4, 0, 0, 0, 0, 133, 0}, {5, 0, 1, 0, 0, 0, 0, 0, 0}, {5, 0, 1, 0, 0, 0, 0, 133, 1}} {
		if _, err := DecodeResponsePeerRequest(bad); err == nil {
			t.Fatalf("invalid tail admitted: %x", bad)
		}
	}
}

func TestTradeInvitationLeaseAndOverlappingPrompts(t *testing.T) {
	m := NewRuntimePartyManager()
	if _, ok := m.RecordOnlinePeerRequest(5, 1, 1, 1, 1); !ok {
		t.Fatal("initial invite")
	}
	for _, pair := range [][2]uint16{{5, 9}, {9, 1}, {1, 5}, {1, 9}} {
		if _, ok := m.RecordOnlinePeerRequest(pair[0], 1, pair[1], 1, 1); ok {
			t.Fatalf("overlap admitted %v", pair)
		}
	}
	if _, ok := m.ConsumePeerRequest(5, 2, 1, 1, 1); ok {
		t.Fatal("stale generation consumed invitation")
	}
	if _, ok := m.ConsumePeerRequest(5, 1, 1, 1, 1); !ok {
		t.Fatal("valid invitation lost")
	}
	tradeKey := runtimePartyInviteKey{invitee: 5, inviter: 1, mode: 1}
	m.RecordOnlinePeerRequest(5, 1, 1, 1, 1)
	invite := m.invites[tradeKey]
	invite.expiresAt = time.Now().Add(-time.Second)
	m.invites[tradeKey] = invite
	if _, ok := m.ConsumePeerRequest(5, 1, 1, 1, 1); ok {
		t.Fatal("expired invitation accepted")
	}
	if len(m.invites) != 0 {
		t.Fatal("expired invitation retained")
	}
	m.RecordOnlinePeerRequest(5, 1, 1, 1, 1)
	invite = m.invites[tradeKey]
	invite.expiresAt = time.Now().Add(-time.Second)
	m.invites[tradeKey] = invite
	if _, ok := m.RecordOnlinePeerRequest(9, 1, 1, 1, 1); !ok {
		t.Fatal("expired invite blocked fresh request")
	}
	if _, ok := m.invites[tradeKey]; ok {
		t.Fatal("expired peer key retained")
	}
}

func TestRuntimePartyManagerConcurrentApplicationsDoNotOverwrite(t *testing.T) {
	m := NewRuntimePartyManager()
	created := m.Create(testRuntimePartyMember(101, 1), alignedcmd.PartyState{MaxMembers: 4})
	if !created.OK {
		t.Fatalf("create party failed: %+v", created)
	}
	// Cleric (202) and Berserker (303) both apply to leader (101) before leader accepts either.
	reqCleric, ok := m.RecordOnlinePeerRequest(101, 1, 202, 1, 12)
	if !ok || reqCleric.JoinerID != 202 || reqCleric.AnchorID != 101 || reqCleric.PartyID != created.Party.ID {
		t.Fatalf("cleric application failed: ok=%t req=%+v", ok, reqCleric)
	}
	reqBerserker, ok := m.RecordOnlinePeerRequest(101, 1, 303, 1, 12)
	if !ok || reqBerserker.JoinerID != 303 || reqBerserker.AnchorID != 101 || reqBerserker.PartyID != created.Party.ID {
		t.Fatalf("berserker application failed: ok=%t req=%+v", ok, reqBerserker)
	}
	// Leader accepts Berserker first, then Cleric — both must succeed!
	if got, ok := m.ConsumePeerRequest(101, 1, 303, 1, 12); !ok || got.JoinerID != 303 {
		t.Fatalf("consume berserker failed: ok=%t req=%+v", ok, got)
	}
	if got, ok := m.ConsumePeerRequest(101, 1, 202, 1, 12); !ok || got.JoinerID != 202 {
		t.Fatalf("consume cleric failed (overwritten by berserker!): ok=%t req=%+v", ok, got)
	}
}
