package mysql

// 说明：本文件原属上游 #87（邮件/邮件箱增量写）。该批含组队文件，按移植规则未整批引入；
// 这里只保留被其它测试复用的公共 helper，依赖 #87 生产实现的两个用例不移植。

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"longheng.io/server/internal/modules/dnf/repository"
)

// These tests drive the production WithinMailboxAssets transaction through the
// row model connector already used by the account item delta suite. They only
// touch the single assetCharacterID inventory that the unit of work exposes;
// the mailbox character's inventory rows are seeded to prove they stay
// untouched.
const (
	mailboxDeltaAssetOwner   = "77"
	mailboxDeltaMailboxOwner = "88"
)

func mailboxDeltaUnitOfWork(state *cargoDeltaSQLState) (*mysqlMailboxAssetUnitOfWork, *sql.DB) {
	database := sql.OpenDB(cargoDeltaConnector{state})
	u := &mysqlMailboxAssetUnitOfWork{router: mysqlRouter{
		db:          SQLDBAdapter{DB: database},
		readDBs:     []string{"replica"},
		writeDBs:    []string{"dnf_s1_w1"},
		tablePrefix: "dnf",
		now:         time.Now,
	}}
	return u, database
}

func seedMailboxDeltaInventory(t *testing.T, db *rowModelDB, characterID string, rows int) repository.InventoryRecord {
	t.Helper()
	record := repository.InventoryRecord{
		CharacterID: characterID,
		Slots:       map[string]repository.ItemStack{},
		Warehouse:   map[string]repository.ItemStack{},
	}
	for n := 0; n < rows; n++ {
		record.Slots[bagMoveSlotKey(0, int16(n))] = bagMoveStackable(int64(700000+n), 5)
	}
	record.Warehouse[bagMoveSlotKey(2, 0)] = bagMoveEquipStack(900000)
	record.Warehouse[bagMoveSlotKey(2, 1)] = bagMoveStackable(900001, 4)
	if err := newRowModelGroup(t, db, false).Inventory.Save(context.Background(), record); err != nil {
		t.Fatalf("seed %s inventory: %v", characterID, err)
	}
	db.calls = nil // Measure the transaction under test, not fixture creation.
	return record
}

func mailboxDeltaWrites(calls []fakeSQLCall) []fakeSQLCall {
	out := make([]fakeSQLCall, 0, len(calls))
	for _, call := range calls {
		if strings.HasPrefix(call.Query, "SELECT ") {
			continue
		}
		out = append(out, call)
	}
	return out
}

func mailboxDeltaWriteStats(calls []fakeSQLCall) (statements, parameters int) {
	for _, call := range mailboxDeltaWrites(calls) {
		statements++
		parameters += len(call.Args)
	}
	return statements, parameters
}

// mailboxDeltaOwnerLines keeps every snapshot line belonging to one character
// row (inventories, inventory_items and inventory_item_extra). Parent rows have
// no trailing columns; item rows start with the same "character_id=" prefix.
func mailboxDeltaOwnerLines(snapshot map[string][]string, owner string) map[string][]string {
	out := map[string][]string{}
	exact := "character_id=" + owner
	withNext := exact + ","
	for table, lines := range snapshot {
		for _, line := range lines {
			// Parent rows hold only the owner id; item rows sort bind_flag
			// first, so the owner field can appear anywhere in the line.
			if line == exact || strings.Contains(line, withNext) {
				out[table] = append(out[table], line)
			}
		}
	}
	return out
}

func mailboxDeltaWithoutEntry(lines map[string][]string, entryKey string) map[string][]string {
	out := map[string][]string{}
	needle := "entry_key=" + entryKey + ","
	for table, rows := range lines {
		kept := make([]string, 0, len(rows))
		for _, row := range rows {
			if strings.Contains(row, needle) {
				continue
			}
			kept = append(kept, row)
		}
		if len(kept) > 0 {
			out[table] = kept
		}
	}
	return out
}
