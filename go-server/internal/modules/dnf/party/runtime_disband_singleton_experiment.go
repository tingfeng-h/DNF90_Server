package party

// Graduated production policy:
//   - Dungeon 2 -> 1 departures (LeaveDungeonExpected / KickInContext inDungeon=true) disband.
//   - Town 2 -> 1 departures (Leave / KickInContext inDungeon=false) retain the surviving
//     member in place with its fixed slot.
func disbandOnSingleton() bool {
	return true
}

func KeepTownSingletonEnabled() bool {
	return true
}

func disbandOnSingletonTown() bool {
	return false
}
