package ziterate

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

type IPRanges struct {
	Networks []net.IPNet
}

func ReadIPRanges(r io.Reader) (ranges IPRanges, err error) {
	buf := make([]byte, 0)
	scanner := bufio.NewScanner(r)
	scanner.Buffer(buf, 2048)
	out := IPRanges{}
	for scanner.Scan() {
		row := scanner.Text()
		// Strip comments and spaces
		row = strings.Split(row, "#")[0]
		row = strings.TrimSpace(row)
		if strings.Contains(row, "/") {
			_, ipNet, err := net.ParseCIDR(row)
			if err != nil {
				return ranges, err
			}
			out.Networks = append(out.Networks, *ipNet)
		} else {
			ip := net.ParseIP(row)
			if ip == nil {
				return ranges, fmt.Errorf("invalid IP: %s", row)
			}
			mask := net.CIDRMask(8*len(ip), 8*len(ip))
			out.Networks = append(out.Networks, net.IPNet{
				IP:   ip,
				Mask: mask,
			})
		}
	}
	err = scanner.Err()
	return out, err
}
