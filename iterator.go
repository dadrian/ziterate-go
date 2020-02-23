package ziterate

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type BigIntGroupIterator struct {
	g         *Group
	generator *big.Int
	start     *big.Int
	end       *big.Int
	current   *big.Int
}

func groupIteratorFromGroup(g *Group) (*BigIntGroupIterator, error) {
	generator, err := g.findMultiplicativeGenerator()
	if err != nil {
		return nil, err
	}
	start, err := rand.Int(rand.Reader, g.P)
	if err != nil {
		return nil, err
	}
	res := &BigIntGroupIterator{
		g:         g,
		generator: generator,
		start:     big.NewInt(0).Add(big.NewInt(0), start),
		end:       big.NewInt(0).Add(big.NewInt(0), start),
		current:   big.NewInt(0).Add(big.NewInt(0), start),
	}
	return res, nil
}

func (it *BigIntGroupIterator) next() *big.Int {
	if it.current == nil {
		return nil
	}
	it.current.Mul(it.current, it.generator)
	it.current.Mod(it.current, it.g.P)

	out := it.current
	if it.current.Cmp(it.end) == 0 {
		it.current = nil
	}
	return out
}

func (it *BigIntGroupIterator) Next() interface{} {
	out := it.next()
	if out == nil {
		return nil
	}
	return out
}

type UintGroupIterator struct {
	g         *Group
	prime     uint64
	generator uint32
	start     uint64
	end       uint64
	current   uint64
}

const (
	PrimeBoundForSmallGroup   = (1 << 40)
	MaxGeneratorForSmallGroup = (1 << 24)
)

func smallGroupIteratorFromGroup(g *Group) (*UintGroupIterator, error) {
	maxP := big.NewInt(PrimeBoundForSmallGroup)
	if g.P.Cmp(maxP) > 0 || !g.P.IsUint64() {
		return nil, fmt.Errorf("prime %s is too big", g.P)
	}
	p := g.P.Uint64()
	var generator uint32
	maxGenerator := big.NewInt(MaxGeneratorForSmallGroup)
	for {
		gen, err := g.findMultiplicativeGenerator()
		if err != nil {
			return nil, err
		}
		if gen.Cmp(maxGenerator) == 1 {
			continue
		}
		generator = uint32(gen.Uint64())
		break
	}
	start, err := rand.Int(rand.Reader, g.P)
	if err != nil {
		return nil, err
	}
	return &UintGroupIterator{
		g:         g,
		prime:     p,
		generator: generator,
		start:     start.Uint64(),
		end:       start.Uint64(),
		current:   start.Uint64(),
	}, nil
}

func (it *UintGroupIterator) next() uint64 {
	if it.current == 0 {
		return 0
	}
	it.current *= uint64(it.generator)
	it.current %= it.prime
	out := it.current
	if it.current == it.end {
		it.current = 0
	}
	return out
}

func (it *UintGroupIterator) Next() interface{} {
	out := it.next()
	if out == 0 {
		return nil
	}
	return out
}
