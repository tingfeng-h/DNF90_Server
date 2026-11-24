package dnfbridge

// v15 (Scheme 2) supersedes v14's disband-then-recreate town singleton retention.
// Town 2 -> 1 retention is now handled in place by internal/modules/dnf/party
// (KeepTownSingletonEnabled / disbandOnSingletonTown / Leave / KickInContext),
// preserving fixed slots and avoiding false "队伍已解散" disband notices,
// while dungeon 2 -> 1 (LeaveDungeonExpected and KickInContext with inDungeon=true)
// continues to disband.

// dungeonRuntime 是 runtimeDungeonState 的可读别名。
type dungeonRuntime = runtimeDungeonState
