package ziterate

import (
	"fmt"
	"math/big"
)

// ZMapGroups contains the cyclic groups used by ZMap, ordered by prime size.
var ZMapGroups = []*Group{
	{
		// 2^8 + 1
		P:            big.NewInt(257),
		KnownRoot:    big.NewInt(3),
		OrderFactors: []*big.Int{big.NewInt(2)},
	},
	{
		// 2^16 + 1
		P:            big.NewInt(65537),
		KnownRoot:    big.NewInt(3),
		OrderFactors: []*big.Int{big.NewInt(2)},
	},
	{
		// 2^24 + 43
		P:            big.NewInt(16777259),
		KnownRoot:    big.NewInt(2),
		OrderFactors: []*big.Int{big.NewInt(2), big.NewInt(23), big.NewInt(103), big.NewInt(3541)},
	},
	{
		// 2^28 + 3
		P:            big.NewInt(268435459),
		KnownRoot:    big.NewInt(2),
		OrderFactors: []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(19), big.NewInt(87211)},
	},
	{
		// 2^32 + 15
		P:            big.NewInt(4294967311),
		KnownRoot:    big.NewInt(3),
		OrderFactors: []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(5), big.NewInt(131), big.NewInt(364289)},
	},
	{
		// 2^33 + 17
		P:            big.NewInt(8589934609),
		KnownRoot:    big.NewInt(19),
		OrderFactors: []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(59), big.NewInt(3033169)},
	},
	{
		// 2^34 + 25
		P:            big.NewInt(17179869209),
		KnownRoot:    big.NewInt(3),
		OrderFactors: []*big.Int{big.NewInt(2), big.NewInt(83), big.NewInt(1277), big.NewInt(20261)},
	},
	{
		// 2^36 + 31
		P:            big.NewInt(68719476767),
		KnownRoot:    big.NewInt(5),
		OrderFactors: []*big.Int{big.NewInt(2), big.NewInt(163), big.NewInt(883), big.NewInt(238727)},
	},
	{
		// 2^40 + 15
		P:            big.NewInt(1099511627791),
		KnownRoot:    big.NewInt(3),
		OrderFactors: []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(5), big.NewInt(36650387593)},
	},
	{
		// 2^44 + 7
		P:            big.NewInt(17592186044423),
		KnownRoot:    big.NewInt(5),
		OrderFactors: []*big.Int{big.NewInt(2), big.NewInt(11), big.NewInt(53), big.NewInt(97), big.NewInt(155542661)},
	},
	{
		// 2^48 + 23
		P:            big.NewInt(281474976710677),
		KnownRoot:    big.NewInt(6),
		OrderFactors: []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(7), big.NewInt(1361), big.NewInt(2462081249)},
	},
}

func SmallestZMapGroupFor(n uint64) (*Group, error) {
	for _, g := range ZMapGroups {
		if g.P.Uint64() > n {
			return g, nil
		}
	}
	return nil, fmt.Errorf("no ZMap group contains %d elements", n)
}
