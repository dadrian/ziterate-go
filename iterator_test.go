package ziterate

import (
	"crypto/rand"
	"math/big"
	"math/bits"
	"strconv"
	"testing"
)

type elementToStringFunc func(interface{}) string

func largestUintGroup() *Group {
	limit := big.NewInt(PrimeBoundForSmallGroup)
	for i := len(ZMapGroups) - 1; i >= 0; i-- {
		if ZMapGroups[i].P.Cmp(limit) <= 0 {
			return ZMapGroups[i]
		}
	}
	return nil
}

func testIteratorInterface(t *testing.T, it Iterator, size int, elementToString elementToStringFunc) {
	count := 0
	counts := make(map[string]int)
	for i := it.Next(); i != nil; i = it.Next() {
		s := elementToString(i)
		counts[s]++
		count++
		if count > size {
			break
		}
	}
	if count != size {
		t.Errorf("expected %d iterations, got %d", size, count)
	}
	if len(counts) != size {
		t.Errorf("expected %d unique elements, got %d", size, len(counts))
	}
	for s, n := range counts {
		if n != 1 {
			t.Errorf("count for %s not 1: got %d", s, n)
		}
	}
}

func TestBigIntIterator(t *testing.T) {
	g := ZMapGroups[0]
	it, err := BigIntGroupIteratorFromGroup(g, rand.Reader)
	t.Log(it.generator)
	if err != nil {
		t.Fatal(err)
	}
	if err := g.checkIfMultiplicativeGenerator(it.generator); err != nil {
		t.Fatal(err)
	}
	toString := func(i interface{}) string {
		return i.(*big.Int).String()
	}
	testIteratorInterface(t, it, 256, toString)
}

func TestSmallGroupIterator(t *testing.T) {
	g := ZMapGroups[0]
	it, err := UintGroupIteratorFromGroup(g, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(it.generator)
	if err := g.checkIfMultiplicativeGenerator(big.NewInt(int64(it.generator))); err != nil {
		t.Fatal(err)
	}
	toString := func(i interface{}) string {
		u := i.(uint64)
		return strconv.FormatUint(u, 10)
	}
	testIteratorInterface(t, it, 256, toString)
}

func TestUintGroupIteratorZMapGroupBounds(t *testing.T) {
	accepted := largestUintGroup()
	if hi, _ := bits.Mul64(accepted.P.Uint64(), MaxGeneratorForSmallGroup); hi != 0 {
		t.Fatalf("expected %s * %d not to overflow uint64", accepted.P, MaxGeneratorForSmallGroup)
	}
	if _, err := UintGroupIteratorFromGroup(accepted, rand.Reader); err != nil {
		t.Fatalf("expected %s to be accepted: %s", accepted.P, err)
	}

	rejected := ZMapGroups[0]
	for _, g := range ZMapGroups {
		if g.P.Cmp(accepted.P) > 0 {
			rejected = g
			break
		}
	}
	if hi, _ := bits.Mul64(rejected.P.Uint64(), MaxGeneratorForSmallGroup); hi == 0 {
		t.Fatalf("expected %s * %d to overflow uint64", rejected.P, MaxGeneratorForSmallGroup)
	}
	if _, err := UintGroupIteratorFromGroup(rejected, rand.Reader); err == nil {
		t.Fatalf("expected %s to be rejected", rejected.P)
	}
}

func BenchmarkIteratorFullBigInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ZMapGroups[0]
		it, err := BigIntGroupIteratorFromGroup(g, rand.Reader)
		if err != nil {
			b.Fatal(err)
		}
		for x := it.NextBigInt(); x != nil; x = it.NextBigInt() {
		}
	}
}

func BenchmarkIteratorFullBigIntInterface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ZMapGroups[0]
		bg, err := BigIntGroupIteratorFromGroup(g, rand.Reader)
		var it Iterator
		it = bg
		if err != nil {
			b.Fatal(err)
		}
		for x := it.Next(); x != nil; x = it.Next() {
		}
	}
}

func BenchmarkIteratorFullUint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ZMapGroups[0]
		it, err := UintGroupIteratorFromGroup(g, rand.Reader)
		if err != nil {
			b.Fatal(err)
		}
		for x := it.NextUint(); x != 0; x = it.NextUint() {
		}
	}
}

func BenchmarkIteratorFullUint64Interface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := ZMapGroups[0]
		sg, err := UintGroupIteratorFromGroup(g, rand.Reader)
		var it Iterator
		it = sg
		if err != nil {
			b.Fatal(err)
		}
		for x := it.Next(); x != nil; x = it.Next() {
		}
	}
}

func BenchmarkIteratorNextBigInt(b *testing.B) {
	g := ZMapGroups[len(ZMapGroups)-1]
	it, err := BigIntGroupIteratorFromGroup(g, rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		x := it.NextBigInt()
		if x == nil {
			b.Fatal("finished before bench")
		}
	}
}

func BenchmarkIteratorNextBigIntInterface(b *testing.B) {
	g := ZMapGroups[len(ZMapGroups)-1]
	bigIt, err := BigIntGroupIteratorFromGroup(g, rand.Reader)
	var it Iterator
	it = bigIt
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		x := it.Next()
		if x == nil {
			b.Fatal("finished before bench")
		}
	}
}

func BenchmarkIteratorNextUint64(b *testing.B) {
	g := largestUintGroup()
	it, err := UintGroupIteratorFromGroup(g, rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		x := it.NextUint()
		if x == 0 {
			b.Fatal("finished before bench")
		}
	}
}

func BenchmarkIteratorNextUint64Interface(b *testing.B) {
	g := largestUintGroup()
	sg, err := UintGroupIteratorFromGroup(g, rand.Reader)
	var it Iterator
	it = sg
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		x := it.Next()
		if x == nil {
			b.Fatal("finished before bench")
		}
	}
}
