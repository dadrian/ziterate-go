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
	validGroups := zmapGroups
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

func TestIsCoprime(t *testing.T) {
	g := Group{
		P:            big.NewInt(23),
		KnownRoot:    big.NewInt(5),
		OrderFactors: []*big.Int{big.NewInt(2), big.NewInt(11)},
	}
	if err := g.IsValid(); err != nil {
		t.Fatalf("invalid group: %s", err)
	}
	coprimes := []int64{
		3, 5, 7, 9, 13, 15, 17, 19, 21,
	}
	shared := []int64{
		2, 4, 6, 8, 10, 11, 12, 14, 16, 18, 20,
	}
	for _, c := range coprimes {
		if !g.isCoprime(big.NewInt(c)) {
			t.Errorf("%d should be coprime with 22", c)
		}
	}
	for _, s := range shared {
		if g.isCoprime(big.NewInt(s)) {
			t.Errorf("%d is not coprime with 22", s)
		}
	}

}

func TestFindAdditiveGenerator(t *testing.T) {
	g := zmapGroups[0]
	additiveGenerator, err := g.findAdditiveGenerator()
	if err != nil {
		t.Fatalf("%s", err)
	}
	if !g.isCoprime(additiveGenerator) {
		t.Errorf("%s should be coprime with p (%s)", additiveGenerator, g.P)
	}
}

func TestFindMultiplicativeGenerator(t *testing.T) {
	g := zmapGroups[0]
	mg, err := g.findMultiplicativeGenerator()
	if err != nil {
		t.Fatal(err)
	}
	if err := g.checkIfMultiplicativeGenerator(mg); err != nil {
		t.Errorf("%s is not a multiplicative generator: %s", mg, err)
	}
}
