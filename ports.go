package ziterate

import (
	"fmt"
	"strconv"
	"strings"
)

// TargetPorts stores parsed target ports. IncludePort controls whether ports
// should be printed; an empty user definition uses one synthetic zero port.
type TargetPorts struct {
	Ports       []uint16
	IncludePort bool
}

// ParseTargetPorts parses ZMap-style port definitions.
func ParseTargetPorts(def string) (TargetPorts, error) {
	def = strings.TrimSpace(def)
	if def == "" {
		return TargetPorts{Ports: []uint16{0}}, nil
	}
	if def == "*" {
		ports := make([]uint16, 0, 1<<16)
		for i := 0; i <= 0xffff; i++ {
			ports = append(ports, uint16(i))
		}
		return TargetPorts{Ports: ports, IncludePort: true}, nil
	}

	var ports []uint16
	for _, part := range strings.Split(def, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			return TargetPorts{}, fmt.Errorf("empty port in %q", def)
		}
		if strings.Contains(part, "-") {
			bounds := strings.Split(part, "-")
			if len(bounds) != 2 {
				return TargetPorts{}, fmt.Errorf("invalid port range: %s", part)
			}
			first, err := parsePort(bounds[0])
			if err != nil {
				return TargetPorts{}, err
			}
			last, err := parsePort(bounds[1])
			if err != nil {
				return TargetPorts{}, err
			}
			if first > last {
				return TargetPorts{}, fmt.Errorf("invalid port range: %d-%d", first, last)
			}
			for port := first; port <= last; port++ {
				ports = append(ports, uint16(port))
			}
			continue
		}
		port, err := parsePort(part)
		if err != nil {
			return TargetPorts{}, err
		}
		ports = append(ports, uint16(port))
	}
	return TargetPorts{Ports: ports, IncludePort: true}, nil
}

func parsePort(s string) (int, error) {
	s = strings.TrimSpace(s)
	port, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid port: %s", s)
	}
	if port < 0 || port > 0xffff {
		return 0, fmt.Errorf("port out of range: %d", port)
	}
	return port, nil
}
