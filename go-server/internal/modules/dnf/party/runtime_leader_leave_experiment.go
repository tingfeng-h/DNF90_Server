package party

// Graduated production policy: when the party leader leaves with surviving
// members, transfer leadership within the same PartyID so surviving members
// keep their fixed slots and UDP relay ports without re-handshaking.
func leaderLeaveSameID() bool {
	return true
}
