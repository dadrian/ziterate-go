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
