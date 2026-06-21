package ziterate

import "testing"

func TestSmallestZMapGroupFor(t *testing.T) {
	tests := []struct {
		name      string
		n         uint64
		wantPrime uint64
	}{
		{
			name:      "zero",
			n:         0,
			wantPrime: 257,
		},
		{
			name:      "fits in first group order",
			n:         256,
			wantPrime: 257,
		},
		{
			name:      "equal to first prime needs next group",
			n:         257,
			wantPrime: 65537,
		},
		{
			name:      "fits in 32-bit group order",
			n:         4294967310,
			wantPrime: 4294967311,
		},
		{
			name:      "equal to 32-bit group prime needs next group",
			n:         4294967311,
			wantPrime: 8589934609,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g, err := SmallestZMapGroupFor(tc.n)
			if err != nil {
				t.Fatal(err)
			}
			if got := g.P.Uint64(); got != tc.wantPrime {
				t.Fatalf("got prime %d, want %d", got, tc.wantPrime)
			}
		})
	}
}

func TestSmallestZMapGroupForTooLarge(t *testing.T) {
	largest := ZMapGroups[len(ZMapGroups)-1]
	if _, err := SmallestZMapGroupFor(largest.P.Uint64()); err == nil {
		t.Fatal("expected error for n equal to largest ZMap prime")
	}
}
