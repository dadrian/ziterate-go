package ziterate

import (
	"net"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gotest.tools/assert"
)

func TestReadIPRanges(t *testing.T) {
	inputs := `1.2.3.4
141.212.118.0/24
192.168.0.0/16
10.0.0.0/8
`
	r := strings.NewReader(inputs)
	ranges, err := ReadIPRanges(r)
	assert.NilError(t, err)
	expectedNets := []net.IPNet{
		net.IPNet{
			IP:   net.ParseIP("1.2.3.4"),
			Mask: net.CIDRMask(32, 32),
		},
		net.IPNet{
			IP:   net.ParseIP("141.212.118.0"),
			Mask: net.CIDRMask(24, 32),
		},
		net.IPNet{
			IP:   net.ParseIP("192.168.0.0"),
			Mask: net.CIDRMask(16, 32),
		},
		net.IPNet{
			IP:   net.ParseIP("10.0.0.0"),
			Mask: net.CIDRMask(8, 32),
		},
	}
	assert.Equal(t, len(ranges.Networks), len(expectedNets))
	for idx := range expectedNets {
		assert.Check(t, cmp.Equal(ranges.Networks[idx].String(), expectedNets[idx].String()))
	}
}
