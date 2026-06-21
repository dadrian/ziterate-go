package ziterate

import (
	"strings"
	"testing"
)

func TestReadIPRanges(t *testing.T) {
	inputs := `# comment-only line
1.2.3.4
141.212.118.0/24
192.168.0.0/16
10.0.0.0/8 # inline comment

`
	r := strings.NewReader(inputs)
	ranges, err := ReadIPRanges(r)
	if err != nil {
		t.Fatal(err)
	}
	contained := []string{
		"1.2.3.4",
		"141.212.118.10",
		"192.168.100.100",
		"10.1.2.3",
	}
	for _, ip := range contained {
		ok, err := ranges.ContainsString(ip)
		if err != nil {
			t.Fatalf("ContainsString(%q): %v", ip, err)
		}
		if !ok {
			t.Errorf("expected %s to be contained", ip)
		}
	}
	notContained := []string{
		"1.2.3.5",
		"141.212.119.10",
		"172.16.0.1",
	}
	for _, ip := range notContained {
		ok, err := ranges.ContainsString(ip)
		if err != nil {
			t.Fatalf("ContainsString(%q): %v", ip, err)
		}
		if ok {
			t.Errorf("expected %s not to be contained", ip)
		}
	}
}
