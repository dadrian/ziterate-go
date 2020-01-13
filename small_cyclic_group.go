package ziterate

import (
	"crypto/rand"
	"math"
	"math/big"
)

type smallGroup struct {
	p            uint64
	knownRoot    uint32
	orderFactors []uint32
}

func (g *smallGroup) isCoprime(x uint64) bool {
	if x == 0 {
		return false
	}
	for idx := range g.orderFactors {
		factor := uint64(g.orderFactors[idx])
		if x == factor {
			return false
		} else if x < factor {
			if factor%x == 0 {
				return false
			}
		} else if x < factor {
			if x%factor == 0 {
				return false
			}
		}
	}
	return true
}

func (g *smallGroup) additiveToMultiplicativeIsomorphism(x uint32) (uint64, error) {
	bigX := big.NewInt(int64(x))
	p := big.NewInt(int64(g.p))
	knownRoot := big.NewInt(int64(g.knownRoot))
	out := big.NewInt(0)
	out.Exp(knownRoot, bigX, p)
	return uint64(out.Int64()), nil
}

func (g *smallGroup) findAdditiveGenerator() (uint32, error) {
	maxCandidate := g.p
	if maxCandidate > math.MaxUint32 {
		maxCandidate = math.MaxUint32
	}
	candidate, err := rand.Int(rand.Reader, big.NewInt(int64(maxCandidate)))
	if err != nil {
		return 0, err
	}
	c := uint64(candidate.Int64())
	for !g.isCoprime(c) {
		c += 1
		c %= g.p
		if c > maxCandidate {
			c = 1
		}
	}
	return uint32(c), nil
}

func (g *smallGroup) findMultiplicativeGenerator() (uint32, error) {
	return 0, nil
}
