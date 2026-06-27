package ziterate

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"net/netip"
	"os"
	"sort"
	"strings"
)

// IPv4Range is an inclusive IPv4 address range in host byte order.
type IPv4Range struct {
	Start  uint32
	End    uint32
	CumEnd uint64
}

// IPv4RangeSet stores sorted, non-overlapping allowed IPv4 ranges.
type IPv4RangeSet struct {
	ranges []IPv4Range
	total  uint64
}

// IPv4RangeSetOptions configures construction of an IPv4RangeSet.
type IPv4RangeSetOptions struct {
	AllowEntries []string
	AllowFiles   []string
	BlockEntries []string
	BlockFiles   []string
}

// NewIPv4RangeSet constructs an IPv4RangeSet from allowlist and blocklist
// entries. If any allowlist source is provided, the set starts empty. Otherwise
// it starts with the full IPv4 address space. Blocklists are then subtracted.
func NewIPv4RangeSet(opts IPv4RangeSetOptions) (*IPv4RangeSet, error) {
	hasAllowlist := len(opts.AllowEntries) > 0 || len(opts.AllowFiles) > 0
	var allowed []IPv4Range
	if !hasAllowlist {
		allowed = []IPv4Range{{Start: 0, End: math.MaxUint32}}
	}

	allowRanges, err := parseRangeSources(opts.AllowEntries, opts.AllowFiles)
	if err != nil {
		return nil, err
	}
	if hasAllowlist {
		allowed = normalizeRanges(allowRanges)
	}

	blockRanges, err := parseRangeSources(opts.BlockEntries, opts.BlockFiles)
	if err != nil {
		return nil, err
	}
	allowed = subtractRanges(allowed, normalizeRanges(blockRanges))
	allowed = subtractRanges(allowed, []IPv4Range{{Start: 0, End: 0}})
	allowed = withCumulativeCounts(normalizeRanges(allowed))

	total := uint64(0)
	if len(allowed) > 0 {
		total = allowed[len(allowed)-1].CumEnd
	}
	return &IPv4RangeSet{ranges: allowed, total: total}, nil
}

// Count returns the number of allowed IPv4 addresses.
func (s *IPv4RangeSet) Count() uint64 {
	if s == nil {
		return 0
	}
	return s.total
}

// Lookup returns the index-th allowed IPv4 address in host byte order.
func (s *IPv4RangeSet) Lookup(index uint64) (uint32, bool) {
	if s == nil || index >= s.total {
		return 0, false
	}
	i := sort.Search(len(s.ranges), func(i int) bool {
		return s.ranges[i].CumEnd > index
	})
	if i == len(s.ranges) {
		return 0, false
	}
	prevCum := uint64(0)
	if i > 0 {
		prevCum = s.ranges[i-1].CumEnd
	}
	return s.ranges[i].Start + uint32(index-prevCum), true
}

// Ranges returns a copy of the allowed IPv4 ranges.
func (s *IPv4RangeSet) Ranges() []IPv4Range {
	if s == nil {
		return nil
	}
	out := make([]IPv4Range, len(s.ranges))
	copy(out, s.ranges)
	return out
}

// Uint32ToIPv4 converts a host-order uint32 IPv4 address to netip.Addr.
func Uint32ToIPv4(ip uint32) netip.Addr {
	return netip.AddrFrom4([4]byte{
		byte(ip >> 24),
		byte(ip >> 16),
		byte(ip >> 8),
		byte(ip),
	})
}

func parseRangeSources(entries, files []string) ([]IPv4Range, error) {
	var out []IPv4Range
	for _, entry := range entries {
		ranges, err := parseRangeLines(strings.NewReader(entry))
		if err != nil {
			return nil, err
		}
		out = append(out, ranges...)
	}
	for _, path := range files {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		ranges, readErr := parseRangeLines(file)
		closeErr := file.Close()
		if readErr != nil {
			return nil, fmt.Errorf("%s: %w", path, readErr)
		}
		if closeErr != nil {
			return nil, closeErr
		}
		out = append(out, ranges...)
	}
	return out, nil
}

func parseRangeLines(r io.Reader) ([]IPv4Range, error) {
	var out []IPv4Range
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		entry := strings.Split(scanner.Text(), "#")[0]
		fields := strings.Fields(entry)
		if len(fields) == 0 {
			continue
		}
		ipRange, err := parseIPv4Range(fields[0])
		if err != nil {
			return nil, err
		}
		out = append(out, ipRange)
	}
	return out, scanner.Err()
}

func parseIPv4Range(entry string) (IPv4Range, error) {
	if strings.Contains(entry, "/") {
		prefix, err := netip.ParsePrefix(entry)
		if err != nil {
			return IPv4Range{}, err
		}
		prefix = prefix.Masked()
		if !prefix.Addr().Is4() {
			return IPv4Range{}, fmt.Errorf("not an IPv4 prefix: %s", entry)
		}
		start := ipv4AddrToUint32(prefix.Addr())
		bits := prefix.Bits()
		size := uint64(1) << uint(32-bits)
		return IPv4Range{
			Start: start,
			End:   uint32(uint64(start) + size - 1),
		}, nil
	}
	addr, err := netip.ParseAddr(entry)
	if err != nil {
		return IPv4Range{}, err
	}
	if !addr.Is4() {
		return IPv4Range{}, fmt.Errorf("not an IPv4 address: %s", entry)
	}
	ip := ipv4AddrToUint32(addr)
	return IPv4Range{Start: ip, End: ip}, nil
}

func ipv4AddrToUint32(addr netip.Addr) uint32 {
	a := addr.As4()
	return uint32(a[0])<<24 | uint32(a[1])<<16 | uint32(a[2])<<8 | uint32(a[3])
}

func normalizeRanges(ranges []IPv4Range) []IPv4Range {
	if len(ranges) == 0 {
		return nil
	}
	out := make([]IPv4Range, 0, len(ranges))
	sort.Slice(ranges, func(i, j int) bool {
		if ranges[i].Start == ranges[j].Start {
			return ranges[i].End < ranges[j].End
		}
		return ranges[i].Start < ranges[j].Start
	})
	for _, r := range ranges {
		if r.End < r.Start {
			continue
		}
		if len(out) == 0 {
			out = append(out, IPv4Range{Start: r.Start, End: r.End})
			continue
		}
		last := &out[len(out)-1]
		if r.Start <= last.End || (last.End < math.MaxUint32 && r.Start == last.End+1) {
			if r.End > last.End {
				last.End = r.End
			}
			continue
		}
		out = append(out, IPv4Range{Start: r.Start, End: r.End})
	}
	return out
}

func subtractRanges(allowed, blocked []IPv4Range) []IPv4Range {
	if len(allowed) == 0 || len(blocked) == 0 {
		return allowed
	}
	blocked = normalizeRanges(blocked)
	out := make([]IPv4Range, 0, len(allowed))
	blockIdx := 0
	for _, allow := range normalizeRanges(allowed) {
		start := allow.Start
		exhausted := false
		for blockIdx < len(blocked) && blocked[blockIdx].End < start {
			blockIdx++
		}
		for i := blockIdx; i < len(blocked) && blocked[i].Start <= allow.End; i++ {
			block := blocked[i]
			if block.Start > start {
				out = append(out, IPv4Range{Start: start, End: block.Start - 1})
			}
			if block.End == math.MaxUint32 {
				exhausted = true
				break
			}
			start = block.End + 1
			if start > allow.End {
				break
			}
		}
		if !exhausted && start <= allow.End {
			out = append(out, IPv4Range{Start: start, End: allow.End})
		}
	}
	return out
}

func withCumulativeCounts(ranges []IPv4Range) []IPv4Range {
	total := uint64(0)
	for i := range ranges {
		total += uint64(ranges[i].End) - uint64(ranges[i].Start) + 1
		ranges[i].CumEnd = total
	}
	return ranges
}
