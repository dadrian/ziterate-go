package ziterate

import (
	"crypto/rand"
	"math/big"
)

type groupIterator struct {
	g         *group
	generator *big.Int
	start     *big.Int
	end       *big.Int
	current   *big.Int
}

func groupIteratorFromGroup(g *group) (*groupIterator, error) {
	generator, err := g.findMultiplicativeGenerator()
	if err != nil {
		return nil, err
	}
	start, err := rand.Int(rand.Reader, g.p)
	if err != nil {
		return nil, err
	}
	res := &groupIterator{
		g:         g,
		generator: generator,
		start:     big.NewInt(0).Add(big.NewInt(0), start),
		end:       big.NewInt(0).Add(big.NewInt(0), start),
		current:   big.NewInt(0).Add(big.NewInt(0), start),
	}
	return res, nil
}

func (it *groupIterator) Next() *big.Int {
	if it.current == nil {
		return nil
	}
	it.current.Mul(it.current, it.generator)
	it.current.Mod(it.current, it.g.p)

	out := it.current
	if it.current == it.end {
		it.current = nil
	}
	return out
}
