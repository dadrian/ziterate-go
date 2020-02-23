package ziterate

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var zero = big.NewInt(0)

// Group represents a cyclic group module P. It can be used for additive or multiplicative groups.
type Group struct {
	// P is a prime number
	P *big.Int

	// KnownRoot is a known generator of the multiplicative group of order P - 1.
	KnownRoot *big.Int

	// OrderFactors are the prime factors of P - 1.
	OrderFactors []*big.Int
}

func (g *Group) isCoprime(x *big.Int) bool {
	var residue big.Int
	if x.Cmp(zero) == 0 {
		return false
	}
	for _, factor := range g.OrderFactors {
		cmp := x.Cmp(factor)
		if cmp > 0 {
			// X is bigger than the factor
			residue.Mod(x, factor)
			if residue.Cmp(zero) == 0 {
				return false
			}
		} else if cmp < 0 {
			// Factor is bigger than X
			residue.Mod(factor, x)
			if residue.Cmp(zero) == 0 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

// Perform the isomorphism from (Z/pZ)+ to (Z/pZ)*
// Given known primitive root of (Z/pZ)* n, with x in (Z/pZ)+, do:
//
//     f(x) = n^x mod p
//
// The isomorphism in the reverse direction is discrete log, and is therefore
// hard.
func (g *Group) additiveToMultiplicativeIsomorphism(x *big.Int) *big.Int {
	out := big.NewInt(0)
	out.Exp(g.KnownRoot, x, g.P)
	return out
}

func (g *Group) findAdditiveGenerator() (*big.Int, error) {
	candidate, err := rand.Int(rand.Reader, g.P)
	if err != nil {
		return nil, err
	}
	for !g.isCoprime(candidate) {
		candidate.Add(candidate, big.NewInt(1))
		candidate.Mod(candidate, g.P)
	}
	return candidate, nil
}

func (g *Group) findMultiplicativeGenerator() (*big.Int, error) {
	additiveGenerator, err := g.findAdditiveGenerator()
	if err != nil {
		return nil, err
	}
	multiplicativeGenerator := g.additiveToMultiplicativeIsomorphism(additiveGenerator)
	return multiplicativeGenerator, nil
}

// Check that the primitive root is a generator of the multiplicative
// group. It is a generator if it is not of the order of any of the
// subgroups.
func (g *Group) checkIfMultiplicativeGenerator(m *big.Int) error {
	order := big.NewInt(0)
	order.Sub(g.P, big.NewInt(1))
	for _, factor := range g.OrderFactors {
		possibleSubgroupOrder := big.NewInt(0)
		possibleSubgroupOrder.Div(order, factor)
		subgroupCheck := big.NewInt(0)
		subgroupCheck.Exp(g.KnownRoot, possibleSubgroupOrder, g.P)
		if subgroupCheck.Cmp(big.NewInt(1)) == 0 {
			return fmt.Errorf("not a generator: (%s does not generates %s, it has order %s)", g.KnownRoot, g.P, possibleSubgroupOrder)
		}
	}
	return nil
}

// IsValid checks that the Group is well-defined: that P is prime, the KnownRoot
// has the correct order, and that the factors are factors of the order.
func (g *Group) IsValid() error {
	// Check that p is prime
	if !g.P.ProbablyPrime(32) {
		return fmt.Errorf("not prime: %s", g.P)
	}
	order := big.NewInt(0)
	order.Sub(g.P, big.NewInt(1))
	// Check that the factors are actually factors. Factors might be repeated,
	// so we don't expect them to multiply together, but they should all be
	// congruent to zero module the order.
	for _, factor := range g.OrderFactors {
		remainder := big.NewInt(0)
		remainder.Mod(order, factor)
		if remainder.Cmp(big.NewInt(0)) != 0 {
			return fmt.Errorf("factor is not a factor of the order: (%s not a factor of %s)", factor, order)
		}
	}
	if err := g.checkIfMultiplicativeGenerator(g.KnownRoot); err != nil {
		return err
	}
	return nil
}
