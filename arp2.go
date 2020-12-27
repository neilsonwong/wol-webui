package main

import (
	"github.com/mostlygeek/arp"
)

// LoadARPTable Returns a map of [ip]mac in a map[string]string format
func LoadARPTable() map[string]string {
	t := arp.Table()
	return t
}
