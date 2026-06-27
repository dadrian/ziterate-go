package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunDeterministicSeed(t *testing.T) {
	args := []string{"-e", "1234", "-n", "3", "10.0.0.0/29"}
	var first bytes.Buffer
	if err := run(args, &first); err != nil {
		t.Fatal(err)
	}
	var second bytes.Buffer
	if err := run(args, &second); err != nil {
		t.Fatal(err)
	}
	if first.String() != second.String() {
		t.Fatalf("seeded output differed:\n%s\n---\n%s", first.String(), second.String())
	}
	lines := nonEmptyLines(first.String())
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3: %q", len(lines), first.String())
	}
}

func TestRunPortsOutput(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"-e", "1", "-n", "2", "-p", "80,81", "10.0.0.1"}, &out); err != nil {
		t.Fatal(err)
	}
	lines := nonEmptyLines(out.String())
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2: %q", len(lines), out.String())
	}
	for _, line := range lines {
		if !strings.Contains(line, ",") {
			t.Fatalf("expected port output, got %q", line)
		}
	}
}

func TestRunShardingRequiresSeed(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"--shards", "2", "--shard", "1", "10.0.0.0/30"}, &out); err == nil {
		t.Fatal("expected sharding without seed to fail")
	}
}

func TestRunLongFlagsAndPercentMaxTargets(t *testing.T) {
	var out bytes.Buffer
	err := run([]string{
		"--seed", "7",
		"--target-ports", "80,81",
		"--max-targets", "50%",
		"10.0.0.1",
		"10.0.0.2",
	}, &out)
	if err != nil {
		t.Fatal(err)
	}
	lines := nonEmptyLines(out.String())
	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2: %q", len(lines), out.String())
	}
}

func TestRunHelp(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"--help"}, &out); err != nil {
		t.Fatal(err)
	}
	help := out.String()
	if help == "" {
		t.Error("empty help")
	}

}

func nonEmptyLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}
