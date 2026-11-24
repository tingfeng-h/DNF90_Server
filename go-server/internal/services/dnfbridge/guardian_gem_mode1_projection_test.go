package dnfbridge

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestCurrentMode1EquipmentGuardianGemRawDFromItem(t *testing.T) {
	var medal currentItemListEntry
	medal.patchCore(int16(currentGuildMedalActorSlot), 100380058, 1)
	binary.LittleEndian.PutUint16(medal.data[101:103], 101)
	binary.LittleEndian.PutUint16(medal.data[103:105], 202)
	binary.LittleEndian.PutUint16(medal.data[105:107], 303)
	binary.LittleEndian.PutUint16(medal.data[107:109], 404)

	got := currentMode1EquipmentGuardianGemRawDFromItem(nil, medal)
	if len(got) != 8 {
		t.Fatalf("lateRawD len = %d, want 8", len(got))
	}
	if !bytes.Equal(got, medal.data[101:109]) {
		t.Fatalf("lateRawD = %x, want %x", got, medal.data[101:109])
	}

	// Non-medal slot must not emit lateRawD.
	var weapon currentItemListEntry
	weapon.patchCore(12, 100001, 1)
	binary.LittleEndian.PutUint16(weapon.data[101:103], 101)
	if got := currentMode1EquipmentGuardianGemRawDFromItem(nil, weapon); len(got) != 0 {
		t.Fatalf("non-medal slot emitted lateRawD=%x", got)
	}

	// Empty sockets on medal must emit nil (unchanged empty raw block).
	var emptyMedal currentItemListEntry
	emptyMedal.patchCore(int16(currentGuildMedalActorSlot), 100380058, 1)
	if got := currentMode1EquipmentGuardianGemRawDFromItem(nil, emptyMedal); len(got) != 0 {
		t.Fatalf("empty medal emitted lateRawD=%x", got)
	}

	// Explicit rollback switch disables fallback.
	t.Setenv(mode1GuardianGemLateRawDEnv, "off")
	if got := currentMode1EquipmentGuardianGemRawDFromItem(nil, medal); len(got) != 0 {
		t.Fatalf("rollback switch did not disable lateRawD: %x", got)
	}
}
