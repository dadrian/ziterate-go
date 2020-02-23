package ziterate

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// BigIntGroupIterator uses a big.Int to Iterate over cyclic groups of arbitrary
// size.
type BigIntGroupIterator struct {
	g         *Group
	generator *big.Int
	start     *big.Int
	end       *big.Int
	current   *big.Int
}

// BigIntGroupIteratorFromGroup constructs a BigIntGroupIterator given any valid
// group.
func BigIntGroupIteratorFromGroup(g *Group) (*BigIntGroupIterator, error) {
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

// NextBigInt is a typed version of Next. If the BigIntGroupIterator is being
// used directly, and not through the Iterator interface, this function should
// be used to iterate.
func (it *BigIntGroupIterator) NextBigInt() *big.Int {
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

// Next implements the Iterator interface.
func (it *BigIntGroupIterator) Next() interface{} {
	out := it.NextBigInt()
	if out == nil {
		return nil
	}
	return out
}

// UintGroupIterator uses a uint64 to iterate through sufficiently-small cyclic
// groups.
type UintGroupIterator struct {
	g         *Group
	prime     uint64
	generator uint32
	start     uint64
	end       uint64
	current   uint64
}

const (
	// PrimeBoundForSmallGroup is the largest P allowed to be used with
	// UintGroupIterator
	PrimeBoundForSmallGroup = (1 << 40)

	// MaxGeneratorForSmallGroup is the largest generator used internally by the
	// UintGroupIterator.
	MaxGeneratorForSmallGroup = (1 << 24)
)

// UintGroupIteratorFromGroup constructs a UintGroupIterator from a Group where
// P is below PrimeBoundForSmallGroup. It is faster when used directly on
// equivalent sized groups.
func UintGroupIteratorFromGroup(g *Group) (*UintGroupIterator, error) {
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

// NextUint is the typed version of next. It is considerably faster than using
// Next.
func (it *UintGroupIterator) NextUint() uint64 {
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

// Next implements the Iterator interface.
func (it *UintGroupIterator) Next() interface{} {
	out := it.NextUint()
	if out == 0 {
		return nil
	}
	return out
}
