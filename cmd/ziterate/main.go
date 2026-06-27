package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"math"
	"math/bits"
	"os"
	"strconv"
	"strings"

	"github.com/zmap/ziterate"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("ziterate", flag.ContinueOnError)
	flags.SetOutput(stdout)

	var blocklistFile string
	flags.StringVar(&blocklistFile, "b", "", "blocklist file")
	flags.StringVar(&blocklistFile, "blocklist-file", "", "blocklist file")
	var allowlistFile string
	flags.StringVar(&allowlistFile, "w", "", "allowlist file")
	flags.StringVar(&allowlistFile, "allowlist-file", "", "allowlist file")
	var portsDef string
	flags.StringVar(&portsDef, "p", "", "target ports")
	flags.StringVar(&portsDef, "target-ports", "", "target ports")
	var seed uint64
	flags.Uint64Var(&seed, "e", 0, "seed")
	flags.Uint64Var(&seed, "seed", 0, "seed")
	var maxTargetsDef string
	flags.StringVar(&maxTargetsDef, "n", "", "max targets")
	flags.StringVar(&maxTargetsDef, "max-targets", "", "max targets")
	var shard uint
	flags.UintVar(&shard, "shard", 0, "shard number")
	var shards uint
	flags.UintVar(&shards, "shards", 1, "total shards")

	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	seedGiven := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "e" || f.Name == "seed" {
			seedGiven = true
		}
	})
	if shards > 1 && !seedGiven {
		return fmt.Errorf("seed is required when sharding")
	}
	if shard >= math.MaxUint16 || shards > math.MaxUint16 {
		return fmt.Errorf("shard values must fit in uint16")
	}

	var randomReader io.Reader = rand.Reader
	if seedGiven {
		randomReader = ziterate.NewSeedReader(seed)
	}

	ports, err := ziterate.ParseTargetPorts(portsDef)
	if err != nil {
		return err
	}

	rangeOpts := ziterate.IPv4RangeSetOptions{
		AllowEntries: flags.Args(),
	}
	if allowlistFile != "" {
		rangeOpts.AllowFiles = []string{allowlistFile}
	}
	if blocklistFile != "" {
		rangeOpts.BlockFiles = []string{blocklistFile}
	}
	allowed, err := ziterate.NewIPv4RangeSet(rangeOpts)
	if err != nil {
		return err
	}

	targetSpace, err := targetSpaceSize(allowed.Count(), len(ports.Ports))
	if err != nil {
		return err
	}
	maxTargets, err := parseMaxTargets(maxTargetsDef, targetSpace)
	if err != nil {
		return err
	}

	it, err := ziterate.NewTargetIterator(ziterate.TargetIteratorOptions{
		Allowed:    allowed,
		Ports:      ports,
		Random:     randomReader,
		Shard:      uint16(shard),
		Shards:     uint16(shards),
		MaxTargets: maxTargets,
	})
	if err != nil {
		return err
	}

	out := bufio.NewWriter(stdout)
	defer out.Flush()
	for target, ok := it.Next(); ok; target, ok = it.Next() {
		ip := ziterate.Uint32ToIPv4(target.IP)
		if target.HasPort {
			fmt.Fprintf(out, "%s,%d\n", ip, target.Port)
		} else {
			fmt.Fprintln(out, ip)
		}
	}
	return nil
}

func targetSpaceSize(addrCount uint64, portCount int) (uint64, error) {
	hi, lo := bits.Mul64(addrCount, uint64(portCount))
	if hi != 0 {
		return 0, fmt.Errorf("target space is too large")
	}
	return lo, nil
}

func parseMaxTargets(def string, targetSpace uint64) (uint64, error) {
	def = strings.TrimSpace(def)
	if def == "" {
		return 0, nil
	}
	if strings.HasSuffix(def, "%") {
		percent, err := strconv.ParseFloat(strings.TrimSuffix(def, "%"), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid max-targets percentage: %s", def)
		}
		if percent < 0 {
			return 0, fmt.Errorf("invalid max-targets percentage: %s", def)
		}
		return uint64(math.Ceil(float64(targetSpace) * percent / 100)), nil
	}
	out, err := strconv.ParseUint(def, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid max-targets: %s", def)
	}
	return out, nil
}
