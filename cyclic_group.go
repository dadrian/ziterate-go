package ziterate

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

var zero = big.NewInt(0)
var maxGenerator = big.NewInt(MaxGeneratorForSmallGroup)

// Group represents a cyclic group module P. It can be used for additive or multiplicative groups.
type Group struct {
	// P is a prime number
	P *big.Int

	// KnownRoot is a known generator of the multiplicative group of order P - 1.
	KnownRoot *big.Int

	// OrderFactors are the prime factors of P - 1.
	OrderFactors []*big.Int
}

func (g *Group) findMultiplicativeGenerator(random io.Reader) (*big.Int, error) {
	limit := big.NewInt(0).Set(g.P)
	if limit.Cmp(maxGenerator) > 0 {
		limit.Set(maxGenerator)
	}
	candidate, err := rand.Int(random, g.P)
	if err != nil {
		return nil, err
	}
	candidate.Add(candidate, big.NewInt(1))
	candidate.Mod(candidate, g.P)
	candidate.Mod(candidate, limit)

	for attempts := big.NewInt(0); attempts.Cmp(limit) < 0; attempts.Add(attempts, big.NewInt(1)) {
		if candidate.Cmp(zero) != 0 && g.checkIfMultiplicativeGenerator(candidate) == nil {
			return big.NewInt(0).Set(candidate), nil
		}
		candidate.Add(candidate, big.NewInt(1))
		candidate.Mod(candidate, g.P)
		candidate.Mod(candidate, limit)
	}
	return nil, fmt.Errorf("could not find multiplicative generator below %s", limit)
}

// Check that the primitive root is a generator of the multiplicative
// group. It is a generator if it is not of the order of any of the
// subgroups.
func (g *Group) checkIfMultiplicativeGenerator(m *big.Int) error {
	if m == nil {
		return fmt.Errorf("not a generator: <nil> is outside [1, %s)", g.P)
	}
	if m.Cmp(zero) <= 0 || m.Cmp(g.P) >= 0 {
		return fmt.Errorf("not a generator: %s is outside [1, %s)", m, g.P)
	}
	order := big.NewInt(0)
	order.Sub(g.P, big.NewInt(1))
	for _, factor := range g.OrderFactors {
		possibleSubgroupOrder := big.NewInt(0)
		possibleSubgroupOrder.Div(order, factor)
		subgroupCheck := big.NewInt(0)
		subgroupCheck.Exp(m, possibleSubgroupOrder, g.P)
		if subgroupCheck.Cmp(big.NewInt(1)) == 0 {
			return fmt.Errorf("not a generator: (%s does not generate %s, it has subgroup order %s)", m, g.P, possibleSubgroupOrder)
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
