package dnfbridge

import (
	"context"
	"os"
	"strconv"
	"strings"

	dnfequip "longheng.io/server/internal/modules/dnf/equip"
	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
)

// C2S op48 `DECREASE_DURABILITY`：客户端在战斗中扣掉某件**穿戴装备**的耐久度后
// 通知服务端。负载是 1 字节，取值 = 装备槽位号（现场高频值 0x0C = 12 = 武器槽，
// 与仓内测试 `mailClaimedEquipWeaponSlot int16 = 12` 一致）。
//
// 改动前这个包**没有处理函数**，于是耐久度只减在客户端显示层，服务端持久值永远
// 不动，下次登录又被服务端值覆盖回去 —— 现场表现就是"打怪装备不掉耐久"。
//
// 本实现只做一件事：把该槽位装备的耐久字减 1（下限 0），同时更新
// `Extra["durability"]` 与 raw entry 里的耐久字（偏移由 equip 包统一管理），
// 然后保存穿戴记录。不发放任何物品、不改强化等级、不改其他字段。
//
// 开关：DNFBRIDGE_EQUIP_DURABILITY_DECREASE
//
//	默认 on；off / 0 / false / no / legacy = 回到改动前（忽略 op48）。
const equipDurabilityDecreaseEnv = "DNFBRIDGE_EQUIP_DURABILITY_DECREASE"

// 现场（2026-10-26）：耐久扣减只写进了库，客户端在副本里看不到变化，要出图回城
// 重新拉取装备行才显示。这里扣完立刻把该槽位的穿戴行用 op14 原生刷新推回去。
//
// 开关：DNFBRIDGE_EQUIP_DURABILITY_REFRESH（默认 on；off = 只写库不刷新）
const equipDurabilityRefreshEnv = "DNFBRIDGE_EQUIP_DURABILITY_REFRESH"

func equipDurabilityRefreshEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(equipDurabilityRefreshEnv))) {
	case "0", "off", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

func equipDurabilityDecreaseEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(equipDurabilityDecreaseEnv))) {
	case "0", "off", "false", "no", "legacy":
		return false
	default:
		return true
	}
}

// handleCurrentDecreaseDurability 处理 op48。
func (s *Service) handleCurrentDecreaseDurability(session *gameSession, body []byte) error {
	if s == nil || session == nil || session.selectedCharacterID == 0 || len(body) != 1 {
		return nil
	}
	if !equipDurabilityDecreaseEnabled() {
		return nil
	}
	slot := int16(body[0])
	repositories, ok := s.repositoryGroup()
	if !ok || repositories.Equipment == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), createWriteTimeout)
	defer cancel()
	characterID := strconv.Itoa(int(session.selectedCharacterID))
	record, found, err := repositories.Equipment.Load(ctx, characterID)
	if err != nil || !found {
		return nil
	}
	key := strconv.Itoa(int(slot))
	entry, ok := record.Entries[key]
	if !ok || entry.ItemID <= 0 {
		s.logGameEvent(session, "game-equip-durability-decrease-skipped",
			"slot", slot, "reason", "no_equipment_in_slot")
		return nil
	}
	current, hasWord := dnfequip.RepairDurabilityWord(entry.RawEntry, entry.ItemID)
	if !hasWord {
		if raw := strings.TrimSpace(entry.Extra["durability"]); raw != "" {
			if parsed, parseErr := strconv.ParseInt(raw, 10, 64); parseErr == nil {
				current, hasWord = parsed, true
			}
		}
	}
	if !hasWord {
		s.logGameEvent(session, "game-equip-durability-decrease-skipped",
			"slot", slot, "item_id", entry.ItemID, "raw_len", len(entry.RawEntry), "reason", "durability_evidence_missing")
		return nil
	}
	if current <= 0 {
		s.logGameEvent(session, "game-equip-durability-decrease-skipped",
			"slot", slot, "item_id", entry.ItemID, "reason", "already_broken")
		return nil
	}
	next := current - 1
	dnfequip.SetRepairDurabilityWord(entry.RawEntry, entry.ItemID, next)
	if entry.Extra == nil {
		entry.Extra = make(map[string]string, 1)
	}
	entry.Extra["durability"] = strconv.FormatInt(next, 10)
	record.Entries[key] = entry
	if err := dnfrepo.SaveEquipmentFields(ctx, repositories.Equipment, record, dnfrepo.EquipmentFieldEntries); err != nil {
		s.logGameEvent(session, "game-equip-durability-decrease-failed",
			"slot", slot, "item_id", entry.ItemID, "error", err)
		return nil
	}
	s.logGameEvent(session, "game-equip-durability-decreased",
		"slot", slot,
		"item_id", entry.ItemID,
		"from", current,
		"to", next,
		"env", equipDurabilityDecreaseEnv)
	if equipDurabilityRefreshEnabled() {
		if err := s.sendWornEquipmentMutationSlots(session, "decrease_durability", []int16{slot}); err != nil {
			s.logGameEvent(session, "game-equip-durability-refresh-failed",
				"slot", slot, "error", err, "env", equipDurabilityRefreshEnv)
		}
	}
	return nil
}
