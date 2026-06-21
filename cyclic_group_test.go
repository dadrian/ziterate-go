package ziterate

import (
	"math/big"
	"testing"
)

var notPrimeGroup = &Group{
	P:            big.NewInt(4294967303),
	KnownRoot:    big.NewInt(3),
	OrderFactors: []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(5), big.NewInt(131), big.NewInt(364289)},
}
var notKnownRootGroup = &Group{
	P:            big.NewInt(4294967311),
	KnownRoot:    big.NewInt(30),
	OrderFactors: []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(5), big.NewInt(131), big.NewInt(364289)},
}
var notFactorsGroup = &Group{
	P:            big.NewInt(4294967311),
	KnownRoot:    big.NewInt(3),
	OrderFactors: []*big.Int{big.NewInt(7), big.NewInt(2), big.NewInt(3), big.NewInt(5), big.NewInt(131), big.NewInt(364289)},
}

func TestIsValid(t *testing.T) {
	validGroups := ZMapGroups
	for idx, g := range validGroups {
		if err := g.IsValid(); err != nil {
			t.Errorf("expected valid group at index %d, got error %s", idx, err)
		}
	}
	invalidGroups := []*Group{notPrimeGroup, notKnownRootGroup, notFactorsGroup}
	for idx, g := range invalidGroups {
		if err := g.IsValid(); err == nil {
			t.Errorf("expected error, got valid for group at index %d", idx)
		}
	}
}

func TestCheckIfMultiplicativeGenerator(t *testing.T) {
	g := Group{
		P:            big.NewInt(23),
		KnownRoot:    big.NewInt(5),
		OrderFactors: []*big.Int{big.NewInt(2), big.NewInt(11)},
	}
	if err := g.IsValid(); err != nil {
		t.Fatalf("invalid group: %s", err)
	}
	generators := []int64{
		5, 7, 10, 11, 14, 15, 17, 19, 20, 21,
	}
	notGenerators := []int64{
		1, 2, 3, 4, 6, 8, 9, 12, 13, 16, 18, 22,
	}
	for _, generator := range generators {
		if err := g.checkIfMultiplicativeGenerator(big.NewInt(generator)); err != nil {
			t.Errorf("%d should be a multiplicative generator: %s", generator, err)
		}
	}
	for _, notGenerator := range notGenerators {
		if err := g.checkIfMultiplicativeGenerator(big.NewInt(notGenerator)); err == nil {
			t.Errorf("%d should not be a multiplicative generator", notGenerator)
		}
	}
}

func TestFindMultiplicativeGenerator(t *testing.T) {
	g := ZMapGroups[0]
	mg, err := g.findMultiplicativeGenerator()
	if err != nil {
		t.Fatal(err)
	}
	if err := g.checkIfMultiplicativeGenerator(mg); err != nil {
		t.Errorf("%s is not a multiplicative generator: %s", mg, err)
	}
}
