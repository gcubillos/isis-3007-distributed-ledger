package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain"
	"github.com/perlin-network/noise"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

func main() {
	t := time.Now()
	genesisBlock := blockchain.Block{}
	genesisBlock = blockchain.Block{0, t.String(), 0, blockchain.CalculateHash(genesisBlock), ""}

	mutex.Lock()
	blockchain.Blockchain = append(blockchain.Blockchain, genesisBlock)
	mutex.Unlock()

	// Create a first node in the network
	alice, err := noise.NewNode()
	check(err)

	bob, err := noise.NewNode()
	check(err)

	bob.Handle(func(ctx noise.HandlerContext) error {
		if !ctx.IsRequest() {
			return nil
		}

		fmt.Printf("Got a message from Alice: '%s'\n", string(ctx.Data()))

		return nil
	})

	alice.Handle(func(ctx noise.HandlerContext) error {
		if !ctx.IsRequest() {
			return nil
		}

		fmt.Printf("Got a message from Bob: '%s'\n", string(ctx.Data()))

		return nil
	})

	err = alice.Listen()
	check(err)
	err = bob.Listen()
	check(err)
	bytes, err := json.Marshal(genesisBlock)
	res, err := alice.Request(context.TODO(), bob.Addr(), bytes)

	fmt.Printf("", res)

	fmt.Printf("", alice.Addr(), alice.ID())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func convertToBlockchainNode(pNode noise.Node) {
	pNode.Handle(func(ctx noise.HandlerContext) error {
		if !ctx.IsRequest() {
			return nil
		}

		fmt.Printf("Got a message from Bob: '%s'\n", string(ctx.Data()))

		return nil
	})
}
