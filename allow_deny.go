package ziterate

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/zmap/go-iptree/iptree"
)

type IPRanges struct {
	Tree *iptree.IPTree
}

func NewIPRanges() IPRanges {
	return IPRanges{Tree: iptree.New()}
}

func (ranges *IPRanges) addByString(ipcidr string) error {
	if ranges.Tree == nil {
		ranges.Tree = iptree.New()
	}
	return ranges.Tree.SetByString(ipcidr, true)
}

func (ranges IPRanges) Contains(ip net.IP) (bool, error) {
	if ranges.Tree == nil {
		return false, nil
	}
	_, found, err := ranges.Tree.Get(ip)
	return found, err
}

func (ranges IPRanges) ContainsString(ip string) (bool, error) {
	if ranges.Tree == nil {
		return false, nil
	}
	_, found, err := ranges.Tree.GetByString(ip)
	return found, err
}

func ReadIPRanges(r io.Reader) (ranges IPRanges, err error) {
	buf := make([]byte, 0)
	scanner := bufio.NewScanner(r)
	scanner.Buffer(buf, 2048)
	out := NewIPRanges()
	for scanner.Scan() {
		row := scanner.Text()
		// Strip comments and spaces
		row = strings.Split(row, "#")[0]
		row = strings.TrimSpace(row)
		if row == "" {
			continue
		}
		if strings.Contains(row, "/") {
			_, ipNet, err := net.ParseCIDR(row)
			if err != nil {
				return ranges, err
			}
			if err := out.addByString(ipNet.String()); err != nil {
				return ranges, err
			}
		} else {
			ip := net.ParseIP(row)
			if ip == nil {
				return ranges, fmt.Errorf("invalid IP: %s", row)
			}
			ipv4 := ip.To4()
			if ipv4 == nil {
				return ranges, fmt.Errorf("invalid IPv4: %s", row)
			}
			if err := out.addByString(ipv4.String() + "/32"); err != nil {
				return ranges, err
			}
		}
	}
	err = scanner.Err()
	return out, err
}
