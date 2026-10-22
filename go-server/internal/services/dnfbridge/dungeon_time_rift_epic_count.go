package dnfbridge

import (
	cryptorand "crypto/rand"
	"math/big"
)

// Normal-abyss cleared-encounter epic count weights:
//
//	1 件 = 40% (4/10)
//	2 件 = 30% (3/10)
//	3 件 = 20% (2/10)
//	4 件 = 10% (1/10)
//
// Kept in a separate file so a later upstream update of
// dungeon_time_rift_loot.go only conflicts at the small reward loop.
var normalAbyssEpicCountWeights = [...]int{1, 1, 1, 1, 2, 2, 2, 3, 3, 4}

const normalAbyssEpicDropMin = 1

// normalAbyssEpicCountForRoll maps a uniform 0..9 draw to 1..4 with the
// intended 40/30/20/10 weights. Exposed for deterministic tests.
func normalAbyssEpicCountForRoll(roll int) int {
	if roll < 0 || roll >= len(normalAbyssEpicCountWeights) {
		return normalAbyssEpicDropMin
	}
	return normalAbyssEpicCountWeights[roll]
}

// normalAbyssEpicRandomCount returns the per-recipient epic count for one
// completed normal abyss encounter.
func normalAbyssEpicRandomCount() int {
	n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(normalAbyssEpicCountWeights))))
	if err != nil {
		return normalAbyssEpicDropMin
	}
	return normalAbyssEpicCountForRoll(int(n.Int64()))
}
