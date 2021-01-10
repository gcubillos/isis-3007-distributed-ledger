package main

type blockchain struct {
	Blocks []network.Block
	State  map[string]int
}
