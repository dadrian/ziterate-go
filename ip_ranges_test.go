package ziterate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIPv4RangeSetParsingLookupAndZeroExclusion(t *testing.T) {
	set, err := NewIPv4RangeSet(IPv4RangeSetOptions{
		AllowEntries: []string{
			"# comment-only\n0.0.0.0\n192.0.2.0/30 # inline comment\n198.51.100.7",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := set.Count(), uint64(5); got != want {
		t.Fatalf("Count() = %d, want %d", got, want)
	}
	tests := []struct {
		index uint64
		want  uint32
	}{
		{0, ipv4AddrToUint32(Uint32ToIPv4(0xc0000200))},
		{1, 0xc0000201},
		{2, 0xc0000202},
		{3, 0xc0000203},
		{4, 0xc6336407},
	}
	for _, tc := range tests {
		got, ok := set.Lookup(tc.index)
		if !ok {
			t.Fatalf("Lookup(%d) returned false", tc.index)
		}
		if got != tc.want {
			t.Fatalf("Lookup(%d) = %#x, want %#x", tc.index, got, tc.want)
		}
	}
	if _, ok := set.Lookup(5); ok {
		t.Fatal("Lookup(5) returned true")
	}
}

func TestIPv4RangeSetDefaultFullAndBlockSubtract(t *testing.T) {
	set, err := NewIPv4RangeSet(IPv4RangeSetOptions{
		BlockEntries: []string{"255.255.255.254/31"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := set.Count(), uint64(1<<32-3); got != want {
		t.Fatalf("Count() = %d, want %d", got, want)
	}
	first, ok := set.Lookup(0)
	if !ok || first != 1 {
		t.Fatalf("first allowed = %d, %v; want 1, true", first, ok)
	}
	last, ok := set.Lookup(set.Count() - 1)
	if !ok || last != 0xfffffffd {
		t.Fatalf("last allowed = %#x, %v; want 0xfffffffd, true", last, ok)
	}
}

func TestIPv4RangeSetMergeAndSubtract(t *testing.T) {
	set, err := NewIPv4RangeSet(IPv4RangeSetOptions{
		AllowEntries: []string{"10.0.0.0/31", "10.0.0.2", "10.0.0.4"},
		BlockEntries: []string{"10.0.0.1", "10.0.0.4"},
	})
	if err != nil {
		t.Fatal(err)
	}
	got := set.Ranges()
	want := []IPv4Range{
		{Start: 0x0a000000, End: 0x0a000000, CumEnd: 1},
		{Start: 0x0a000002, End: 0x0a000002, CumEnd: 2},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d ranges, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("range %d = %#v, want %#v", i, got[i], want[i])
		}
	}
}

func TestIPv4RangeSetFiles(t *testing.T) {
	dir := t.TempDir()
	allowFile := filepath.Join(dir, "allow.txt")
	blockFile := filepath.Join(dir, "block.txt")
	if err := os.WriteFile(allowFile, []byte("203.0.113.0/30\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(blockFile, []byte("203.0.113.2\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	set, err := NewIPv4RangeSet(IPv4RangeSetOptions{
		AllowFiles: []string{allowFile},
		BlockFiles: []string{blockFile},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := set.Count(), uint64(3); got != want {
		t.Fatalf("Count() = %d, want %d", got, want)
	}
}
