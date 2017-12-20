package main

// Struct to keep tab on a specific Blockchain instance.
type version struct {
	BestHeight  int
	AddressFrom int
}

var nodeAddress string
var knownNodes = []string{"172.19.0.2", "172.19.0.3", "172.19.0.4", "172.19.0.5"}
