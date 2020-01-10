package ziterate

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var zero = big.NewInt(0)

type group struct {
	p            *big.Int
	knownRoot    *big.Int
	orderFactors []*big.Int
}

func (g *group) isCoprime(x *big.Int) bool {
	var residue big.Int
	if x.Cmp(zero) == 0 {
		return false
	}
	for _, factor := range g.orderFactors {
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
func (g *group) additiveToMultiplicativeIsomorphism(x *big.Int) *big.Int {
	out := big.NewInt(0)
	out.Exp(g.knownRoot, x, g.p)
	return out
}

func (g *group) findAdditiveGenerator() (*big.Int, error) {
	candidate, err := rand.Int(rand.Reader, g.p)
	if err != nil {
		return nil, err
	}
	for !g.isCoprime(candidate) {
		candidate.Add(candidate, big.NewInt(1))
		candidate.Mod(candidate, g.p)
	}
	return candidate, nil
}

func (g *group) findMultiplicativeGenerator() (*big.Int, error) {
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
func (g *group) checkIfMultiplicativeGenerator(m *big.Int) error {
	order := big.NewInt(0)
	order.Sub(g.p, big.NewInt(1))
	for _, factor := range g.orderFactors {
		possibleSubgroupOrder := big.NewInt(0)
		possibleSubgroupOrder.Div(order, factor)
		subgroupCheck := big.NewInt(0)
		subgroupCheck.Exp(g.knownRoot, possibleSubgroupOrder, g.p)
		if subgroupCheck.Cmp(big.NewInt(1)) == 0 {
			return fmt.Errorf("not a generator: (%s does not generates %s, it has order %s)", g.knownRoot, g.p, possibleSubgroupOrder)
		}
	}
	return nil
}

func (g *group) isValid() error {
	// Check that p is prime
	if !g.p.ProbablyPrime(32) {
		return fmt.Errorf("not prime: %s", g.p)
	}
	order := big.NewInt(0)
	order.Sub(g.p, big.NewInt(1))
	// Check that the factors are actually factors. Factors might be repeated,
	// so we don't expect them to multiply together, but they should all be
	// congruent to zero module the order.
	for _, factor := range g.orderFactors {
		remainder := big.NewInt(0)
		remainder.Mod(order, factor)
		if remainder.Cmp(big.NewInt(0)) != 0 {
			return fmt.Errorf("factor is not a factor of the order: (%s not a factor of %s)", factor, order)
		}
	}
	if err := g.checkIfMultiplicativeGenerator(g.knownRoot); err != nil {
		return err
	}
	return nil
}
