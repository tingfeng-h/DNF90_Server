package dnfbridge

import (
	"encoding/binary"
)

// Graduated production policy: always project a worn guild medal's 4 guardian-gem
// socket words (raw 0x77 offsets 101..109) into mode1/mode8 `lateRawD`.
func mode1GuardianGemLateRawDEnabled() bool {
	return true
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
