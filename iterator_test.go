package ziterate

import (
	"math/big"
	"testing"
)

func TestIterator(t *testing.T) {
	g := zmapGroups[0]
	it, err := groupIteratorFromGroup(g)
	t.Log(it.generator)
	if err != nil {
		t.Fatal(err)
	}
	if err := g.checkIfMultiplicativeGenerator(it.generator); err != nil {
		t.Fatal(err)
	}
	count := 0
	counts := make([]int64, 257)
	counts[0] = 1
	for i := it.next(); i != nil; i = it.next() {
		counts[i.Int64()] += 1
		count += 1
		if count > 256 {
			break
		}
	}
	if count != 256 {
		t.Errorf("expected 256 iterations, got %d", count)
	}
	for idx, count := range counts {
		if count != 1 {
			t.Errorf("count for %d not 1: got %d", idx, count)
		}
	}
}

func TestSmallGroupIterator(t *testing.T) {
	g := zmapGroups[0]
	it, err := smallGroupIteratorFromGroup(g)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(it.generator)
	if err := g.checkIfMultiplicativeGenerator(big.NewInt(int64(it.generator))); err != nil {
		t.Fatal(err)
	}
	count := 0
	counts := make([]int64, 257)
	counts[0] = 1
	for i := it.next(); i != 0; i = it.next() {
		counts[i] += 1
		count += 1
		if count > 256 {
			break
		}
	}
	if count != 256 {
		t.Errorf("expected 256 iterations, got %d", count)
	}
	for idx, count := range counts {
		if count != 1 {
			t.Errorf("count for %d not 1: got %d", idx, count)
		}
	}
}

func BenchmarkIteratorFullBigInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := zmapGroups[0]
		it, err := groupIteratorFromGroup(g)
		if err != nil {
			b.Fatal(err)
		}
		for x := it.next(); x != nil; x = it.next() {
		}
	}
}

func BenchmarkIteratorFullBigIntInterface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := zmapGroups[0]
		bg, err := groupIteratorFromGroup(g)
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
		g := zmapGroups[0]
		it, err := smallGroupIteratorFromGroup(g)
		if err != nil {
			b.Fatal(err)
		}
		for x := it.next(); x != 0; x = it.next() {
		}
	}
}

func BenchmarkIteratorFullUint64Interface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := zmapGroups[0]
		sg, err := smallGroupIteratorFromGroup(g)
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
	g := zmapGroups[len(zmapGroups)-1]
	it, err := groupIteratorFromGroup(g)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		x := it.next()
		if x == nil {
			b.Fatal("finished before bench")
		}
	}
}

func BenchmarkIteratorNextBigIntInterface(b *testing.B) {
	g := zmapGroups[len(zmapGroups)-1]
	bigIt, err := groupIteratorFromGroup(g)
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
	g := zmapGroups[len(zmapGroups)-1]
	it, err := smallGroupIteratorFromGroup(g)
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		x := it.next()
		if x == 0 {
			b.Fatal("finished before bench")
		}
	}
}

func BenchmarkIteratorNextUint64Interface(b *testing.B) {
	g := zmapGroups[len(zmapGroups)-1]
	sg, err := smallGroupIteratorFromGroup(g)
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
