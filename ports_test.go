package ziterate

import "testing"

func TestParseTargetPortsEmpty(t *testing.T) {
	ports, err := ParseTargetPorts("")
	if err != nil {
		t.Fatal(err)
	}
	if ports.IncludePort {
		t.Fatal("empty port definition should not include ports in output")
	}
	if len(ports.Ports) != 1 || ports.Ports[0] != 0 {
		t.Fatalf("ports = %#v, want synthetic zero port", ports.Ports)
	}
}

func TestParseTargetPortsListAndRanges(t *testing.T) {
	ports, err := ParseTargetPorts("80,443,100-102")
	if err != nil {
		t.Fatal(err)
	}
	want := []uint16{80, 443, 100, 101, 102}
	if !ports.IncludePort {
		t.Fatal("expected IncludePort")
	}
	if len(ports.Ports) != len(want) {
		t.Fatalf("got %d ports, want %d", len(ports.Ports), len(want))
	}
	for i := range want {
		if ports.Ports[i] != want[i] {
			t.Fatalf("port %d = %d, want %d", i, ports.Ports[i], want[i])
		}
	}
}

func TestParseTargetPortsWildcard(t *testing.T) {
	ports, err := ParseTargetPorts("*")
	if err != nil {
		t.Fatal(err)
	}
	if !ports.IncludePort {
		t.Fatal("expected IncludePort")
	}
	if len(ports.Ports) != 1<<16 {
		t.Fatalf("got %d ports, want %d", len(ports.Ports), 1<<16)
	}
	if ports.Ports[0] != 0 || ports.Ports[len(ports.Ports)-1] != 0xffff {
		t.Fatalf("wildcard endpoints = %d, %d", ports.Ports[0], ports.Ports[len(ports.Ports)-1])
	}
}

func TestParseTargetPortsInvalid(t *testing.T) {
	for _, input := range []string{"105-100", "65536", "-1", "80,,81", "abc"} {
		if _, err := ParseTargetPorts(input); err == nil {
			t.Fatalf("ParseTargetPorts(%q) succeeded", input)
		}
	}
}
