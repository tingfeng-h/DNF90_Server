package dnfbridge

import (
	"encoding/binary"
	"os"
	"strings"
)

const mode1GuardianGemLateRawDEnv = "DNFBRIDGE_MODE1_GUARDIAN_GEM_LATE_RAW_D"

// mode1GuardianGemLateRawDEnabled controls projecting a worn guild medal's
// 4 guardian-gem socket words (raw 0x77 offsets 101..109) into mode1/mode8
// equipment create row `lateRawD` (V212 0x0202A917 -> 0x0243599F/0x024359A8
// writing state+0xAF..0xB6). Enabled by default; set off/0/false/legacy to
// roll back.
func mode1GuardianGemLateRawDEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(mode1GuardianGemLateRawDEnv))) {
	case "0", "off", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

// currentMode1EquipmentGuardianGemRawDFromItem projects the 8-byte guardian gem
// socket block (entry.data[101:109]) for actor slot 32 (guild medal) into
// mode1/mode8 `lateRawD`, so both local scene rebuilds and remote inspect /
// mode8 actor profiles display all 4 inlaid guardian gems without needing an
// op14 replay or unequip/re-equip cycle.
func currentMode1EquipmentGuardianGemRawDFromItem(extra map[string]string, entry currentItemListEntry) []byte {
	if explicit := currentMode1EquipmentExtraHex(extra, currentMode1EquipmentRawBlockMaxLen, "mode1_raw_d_hex", "mode1_state_175_179_raw_hex"); len(explicit) != 0 {
		return explicit
	}
	if !mode1GuardianGemLateRawDEnabled() {
		return nil
	}
	if binary.LittleEndian.Uint16(entry.data[:2]) != uint16(currentGuildMedalActorSlot) {
		return nil
	}
	if !currentGuardianGemRawSocketOccupied(entry.data[:]) {
		return nil
	}
	end := currentGuardianGemRawSocketOffset + currentGuardianGemSocketCount*currentGuardianGemRawSocketWidth
	return append([]byte(nil), entry.data[currentGuardianGemRawSocketOffset:end]...)
}
