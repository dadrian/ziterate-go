package ziterate

import "math/big"

var zmapGroups = []*Group{
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
}
