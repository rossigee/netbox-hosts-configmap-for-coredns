package main

import (
	"testing"
)

func TestGenerateHostsFile(t *testing.T) {
	hosts := []HostEntry{
		{
			IP:          "192.168.1.1",
			Hostname:    "example.com",
			Description: "Test host",
		},
		{
			IP:          "192.168.1.2",
			Hostname:    "example2.com",
			Description: "Test host 2",
		},
	}

	expected := "192.168.1.1 example.com # Test host\n192.168.1.2 example2.com # Test host 2"
	actual := generateHostsFile(hosts)

	if actual != expected {
		t.Errorf("Expected \"%s\", but got \"%s\"", expected, actual)
	}
}
