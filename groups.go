package ziterate

import "math/big"

var zmapGroups = []*group{
	{
		// 2^8 + 1
		p:            big.NewInt(257),
		knownRoot:    big.NewInt(3),
		orderFactors: []*big.Int{big.NewInt(2)},
	},
	{
		// 2^16 + 1
		p:            big.NewInt(65537),
		knownRoot:    big.NewInt(3),
		orderFactors: []*big.Int{big.NewInt(2)},
	},
	{
		// 2^24 + 43
		p:            big.NewInt(16777259),
		knownRoot:    big.NewInt(2),
		orderFactors: []*big.Int{big.NewInt(2), big.NewInt(23), big.NewInt(103), big.NewInt(3541)},
	},
	{
		// 2^28 + 3
		p:            big.NewInt(268435459),
		knownRoot:    big.NewInt(2),
		orderFactors: []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(19), big.NewInt(87211)},
	},
	{
		// 2^32 + 15
		p:            big.NewInt(4294967311),
		knownRoot:    big.NewInt(3),
		orderFactors: []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(5), big.NewInt(131), big.NewInt(364289)},
	},
}
