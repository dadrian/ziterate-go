package ziterate

import "testing"

type sequenceIterator struct {
	values []uint64
	index  int
}

func (it *sequenceIterator) Next() interface{} {
	if it.index >= len(it.values) {
		return nil
	}
	out := it.values[it.index]
	it.index++
	return out
}

func TestTargetIteratorMappingSkipsOutOfRange(t *testing.T) {
	allowed, err := NewIPv4RangeSet(IPv4RangeSetOptions{
		AllowEntries: []string{"10.0.0.1/32", "10.0.0.2/32"},
	})
	if err != nil {
		t.Fatal(err)
	}
	it := &TargetIterator{
		allowed:     allowed,
		ports:       TargetPorts{Ports: []uint16{80, 443}, IncludePort: true},
		iterator:    &sequenceIterator{values: []uint64{5, 1, 2, 3, 4}},
		targetSpace: 4,
		shards:      1,
	}

	want := []Target{
		{IP: 0x0a000001, Port: 80, HasPort: true},
		{IP: 0x0a000001, Port: 443, HasPort: true},
		{IP: 0x0a000002, Port: 80, HasPort: true},
		{IP: 0x0a000002, Port: 443, HasPort: true},
	}
	for i, expected := range want {
		got, ok := it.Next()
		if !ok {
			t.Fatalf("Next(%d) returned false", i)
		}
		if got != expected {
			t.Fatalf("Next(%d) = %#v, want %#v", i, got, expected)
		}
	}
	if _, ok := it.Next(); ok {
		t.Fatal("Next() returned true after sequence exhausted")
	}
}

func TestTargetIteratorShardingAndMaxTargets(t *testing.T) {
	allowed, err := NewIPv4RangeSet(IPv4RangeSetOptions{
		AllowEntries: []string{"10.0.0.1/32", "10.0.0.2/32", "10.0.0.3/32", "10.0.0.4/32"},
	})
	if err != nil {
		t.Fatal(err)
	}
	it := &TargetIterator{
		allowed:     allowed,
		ports:       TargetPorts{Ports: []uint16{0}},
		iterator:    &sequenceIterator{values: []uint64{1, 2, 3, 4}},
		targetSpace: 4,
		shard:       1,
		shards:      2,
		maxTargets:  1,
	}
	got, ok := it.Next()
	if !ok {
		t.Fatal("Next() returned false")
	}
	if got.IP != 0x0a000002 || got.HasPort {
		t.Fatalf("target = %#v, want 10.0.0.2 without port", got)
	}
	if _, ok := it.Next(); ok {
		t.Fatal("Next() returned true after maxTargets")
	}
}
