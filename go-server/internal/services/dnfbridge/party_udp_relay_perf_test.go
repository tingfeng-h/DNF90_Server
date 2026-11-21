package dnfbridge

import "testing"

func TestPartyLANSamePrivateSubnet24(t *testing.T) {
	cases := []struct {
		a, b [4]byte
		want bool
	}{
		{[4]byte{192, 168, 1, 10}, [4]byte{192, 168, 1, 22}, true},
		{[4]byte{10, 0, 5, 2}, [4]byte{10, 0, 5, 99}, true},
		{[4]byte{172, 16, 3, 4}, [4]byte{172, 16, 3, 8}, true},
		{[4]byte{192, 168, 1, 10}, [4]byte{192, 168, 2, 10}, false},
		{[4]byte{8, 8, 8, 8}, [4]byte{8, 8, 8, 9}, false},
		{[4]byte{172, 32, 1, 1}, [4]byte{172, 32, 1, 2}, false},
	}
	for _, tc := range cases {
		if got := partyLANSamePrivateSubnet24(tc.a, tc.b); got != tc.want {
			t.Fatalf("partyLANSamePrivateSubnet24(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}
