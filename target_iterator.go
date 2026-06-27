package ziterate

import (
	"fmt"
	"io"
	"math"
	"math/big"
	"math/bits"
)

// Target is an IPv4 target and optional destination port.
type Target struct {
	IP      uint32
	Port    uint16
	HasPort bool
}

// TargetIteratorOptions configures a TargetIterator.
type TargetIteratorOptions struct {
	Allowed    *IPv4RangeSet
	Ports      TargetPorts
	Random     io.Reader
	Shard      uint16
	Shards     uint16
	MaxTargets uint64
}

// TargetIterator maps cyclic group elements into allowed IPv4 targets.
type TargetIterator struct {
	allowed     *IPv4RangeSet
	ports       TargetPorts
	iterator    Iterator
	targetSpace uint64
	shard       uint16
	shards      uint16
	seen        uint64
	emitted     uint64
	maxTargets  uint64
}

// NewTargetIterator constructs a TargetIterator over the configured allowed
// addresses and ports.
func NewTargetIterator(opts TargetIteratorOptions) (*TargetIterator, error) {
	if opts.Allowed == nil || opts.Allowed.Count() == 0 {
		return nil, fmt.Errorf("no allowed targets")
	}
	if len(opts.Ports.Ports) == 0 {
		opts.Ports = TargetPorts{Ports: []uint16{0}}
	}
	if opts.Shards == 0 {
		opts.Shards = 1
	}
	if opts.Shard >= opts.Shards {
		return nil, fmt.Errorf("shard %d must be less than shards %d", opts.Shard, opts.Shards)
	}
	hi, lo := bits.Mul64(opts.Allowed.Count(), uint64(len(opts.Ports.Ports)))
	if hi != 0 || lo > math.MaxUint64-1 {
		return nil, fmt.Errorf("target space is too large")
	}
	group, err := SmallestZMapGroupFor(lo)
	if err != nil {
		return nil, err
	}
	var it Iterator
	if group.P.Cmp(big.NewInt(PrimeBoundForSmallGroup)) <= 0 {
		it, err = UintGroupIteratorFromGroup(group, opts.Random)
	} else {
		it, err = BigIntGroupIteratorFromGroup(group, opts.Random)
	}
	if err != nil {
		return nil, err
	}
	return &TargetIterator{
		allowed:     opts.Allowed,
		ports:       opts.Ports,
		iterator:    it,
		targetSpace: lo,
		shard:       opts.Shard,
		shards:      opts.Shards,
		maxTargets:  opts.MaxTargets,
	}, nil
}

// Next returns the next target, or false when iteration is complete.
func (it *TargetIterator) Next() (Target, bool) {
	if it.maxTargets > 0 && it.emitted >= it.maxTargets {
		return Target{}, false
	}
	for {
		value, ok := it.nextValue()
		if !ok {
			return Target{}, false
		}
		if value == 0 {
			continue
		}
		index := value - 1
		if index >= it.targetSpace {
			continue
		}
		ipIndex := index / uint64(len(it.ports.Ports))
		portIndex := index % uint64(len(it.ports.Ports))
		ip, ok := it.allowed.Lookup(ipIndex)
		if !ok {
			continue
		}
		seen := it.seen
		it.seen++
		if seen%uint64(it.shards) != uint64(it.shard) {
			continue
		}
		it.emitted++
		return Target{
			IP:      ip,
			Port:    it.ports.Ports[portIndex],
			HasPort: it.ports.IncludePort,
		}, true
	}
}

func (it *TargetIterator) nextValue() (uint64, bool) {
	switch v := it.iterator.(type) {
	case *UintGroupIterator:
		out := v.NextUint()
		return out, out != 0
	case *BigIntGroupIterator:
		out := v.NextBigInt()
		if out == nil || !out.IsUint64() {
			return 0, false
		}
		return out.Uint64(), true
	default:
		out := it.iterator.Next()
		if out == nil {
			return 0, false
		}
		switch typed := out.(type) {
		case uint64:
			return typed, true
		case *big.Int:
			if !typed.IsUint64() {
				return 0, false
			}
			return typed.Uint64(), true
		default:
			return 0, false
		}
	}
}
