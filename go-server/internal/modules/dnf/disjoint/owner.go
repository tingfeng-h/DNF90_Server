package disjoint

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	dnfrepo "longheng.io/server/internal/modules/dnf/repository"
)

var (
	ErrOwnerUnavailable  = errors.New("disjoint owner is unavailable")
	ErrAccountRequired   = errors.New("disjoint account id is required")
	ErrCharacterRequired = errors.New("disjoint character id is required")
	ErrCharacterNotFound = errors.New("disjoint character is not found")
	ErrAccountMismatch   = errors.New("disjoint character account mismatch")
	ErrInventoryNotFound = errors.New("disjoint inventory is not found")
	ErrProjectorRequired = errors.New("disjoint mutation projector is required")
)

type Assets struct {
	AccountInventory *dnfrepo.AccountInventoryRecord
	Inventory        *dnfrepo.InventoryRecord
	Character        *dnfrepo.CharacterRecord
}

type Changes struct {
	AccountInventory bool
	Inventory        bool
	CharacterStats   bool
}

type Projector func(*Assets) (Changes, error)

type Command struct {
	AccountID   string
	CharacterID string
	UpdatedAt   time.Time
	Project     Projector
}

// Owner owns character/account authorization plus the atomic account
// warehouse and character inventory write used by all three gameplays.
type Owner struct {
	settings dnfrepo.SettingsRepository
	assets   dnfrepo.AccountCharacterAssetUnitOfWork
}

func NewOwner(repositories dnfrepo.Group) (*Owner, error) {
	if repositories.AccountAssets == nil {
		return nil, ErrOwnerUnavailable
	}
	return &Owner{settings: repositories.Settings, assets: repositories.AccountAssets}, nil
}

func (o *Owner) DisjointEquipment(ctx context.Context, command Command) error {
	return o.mutate(ctx, "equipment-disjoint", command)
}

func (o *Owner) DisjointAvatar(ctx context.Context, command Command) error {
	return o.mutate(ctx, "avatar-disjoint", command)
}

func (o *Owner) CompoundEmblems(ctx context.Context, command Command) error {
	return o.mutate(ctx, "emblem-compound", command)
}

func (o *Owner) mutate(ctx context.Context, operation string, command Command) error {
	if o == nil || o.assets == nil {
		return ErrOwnerUnavailable
	}
	accountID := strings.TrimSpace(command.AccountID)
	characterID := strings.TrimSpace(command.CharacterID)
	if accountID == "" {
		return ErrAccountRequired
	}
	if characterID == "" {
		return ErrCharacterRequired
	}
	if command.Project == nil {
		return ErrProjectorRequired
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	now := command.UpdatedAt.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return o.assets.WithinAccountCharacterAssets(ctx, accountID, characterID, func(
		accounts dnfrepo.AccountInventoryRepository,
		characters dnfrepo.CharacterRepository,
		inventories dnfrepo.InventoryRepository,
		_ dnfrepo.EquipmentRepository,
	) error {
		if accounts == nil || characters == nil || inventories == nil {
			return ErrOwnerUnavailable
		}
		character, found, err := characters.Load(ctx, characterID)
		if err != nil {
			return err
		}
		if !found || strings.TrimSpace(character.CharacterID) != characterID {
			return fmt.Errorf("%w: character=%s operation=%s", ErrCharacterNotFound, characterID, operation)
		}
		if strings.TrimSpace(character.AccountID) != accountID {
			return fmt.Errorf("%w: account=%s character=%s operation=%s", ErrAccountMismatch, accountID, characterID, operation)
		}
		character = dnfrepo.CloneCharacter(character)
		inventory, found, err := inventories.Load(ctx, characterID)
		if err != nil {
			return err
		}
		if !found || strings.TrimSpace(inventory.CharacterID) != characterID || inventory.Slots == nil {
			return fmt.Errorf("%w: character=%s operation=%s", ErrInventoryNotFound, characterID, operation)
		}
		accountInventory, _, err := accounts.Load(ctx, accountID)
		if err != nil {
			return err
		}
		accountInventory = dnfrepo.CloneAccountInventory(accountInventory)
		accountInventory.AccountID = accountID
		if accountInventory.Slots == nil {
			accountInventory.Slots = make(map[string]dnfrepo.ItemStack)
		}
		if err := dnfrepo.BindInventoryContainerState(ctx, o.settings, &inventory); err != nil {
			return err
		}
		inventory = dnfrepo.CloneInventory(inventory)
		beforeInventory := dnfrepo.CloneInventory(inventory)

		changes, err := command.Project(&Assets{
			AccountInventory: &accountInventory,
			Inventory:        &inventory,
			Character:        &character,
		})
		if err != nil {
			return err
		}
		if changes.CharacterStats {
			character.UpdatedAt = now
			if err := dnfrepo.SaveCharacterFields(ctx, characters, character, dnfrepo.CharacterFieldStats); err != nil {
				return err
			}
		}
		if changes.Inventory {
			inventory.UpdatedAt = now
			// Default scope (collection) still does the exact legacy full rewrite.
			// With DNFBRIDGE_INVENTORY_SAVE_SCOPE=slot, only changed slots/extras
			// are written; this avoids rewriting all 1000+ inventory extra rows
			// for one disassemble/compound operation and removes the long row-lock
			// transactions that block the character SELECT ... FOR UPDATE.
			keys := changedInventorySlotKeys(beforeInventory.Slots, inventory.Slots)
			if err := dnfrepo.SaveInventoryFieldsScoped(
				ctx,
				inventories,
				inventory,
				[]dnfrepo.InventoryField{dnfrepo.InventoryFieldSlots},
				dnfrepo.InventorySlotDeltaKeys{string(dnfrepo.InventoryFieldSlots): keys},
			); err != nil {
				return err
			}
		}
		if changes.AccountInventory {
			accountInventory.UpdatedAt = now
			if err := accounts.Save(ctx, accountInventory); err != nil {
				return err
			}
		}
		return nil
	})
}

// changedInventorySlotKeys returns the union of slots added, removed or changed
// by the projector. The result feeds SaveInventoryFieldsScoped; the default
// collection scope ignores it and keeps the legacy full rewrite.
func changedInventorySlotKeys(before, after map[string]dnfrepo.ItemStack) []string {
	keys := make([]string, 0)
	for key, oldItem := range before {
		newItem, exists := after[key]
		if !exists || !reflect.DeepEqual(oldItem, newItem) {
			keys = append(keys, key)
		}
	}
	for key := range after {
		if _, exists := before[key]; !exists {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	return keys
}
