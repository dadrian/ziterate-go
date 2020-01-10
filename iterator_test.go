package ziterate

import (
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
	for i := it.Next(); i != nil; i = it.Next() {
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

func BenchmarkIteratorFull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := zmapGroups[0]
		it, err := groupIteratorFromGroup(g)
		if err != nil {
			b.Fatal(err)
		}
		for x := it.Next(); x != nil; x = it.Next() {
		}
	}
}

func BenchmarkIteratorNext(b *testing.B) {
	g := zmapGroups[len(zmapGroups)-1]
	it, err := groupIteratorFromGroup(g)
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
