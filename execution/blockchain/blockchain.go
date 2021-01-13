package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain"
	ghost "github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain-ghost"
	"github.com/perlin-network/noise"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

func main() {
	t := time.Now()
	genesisBlock := blockchain.Block{}
	genesisBlock = blockchain.Block{t, blockchain.CalculateHash(genesisBlock), "", 0, make([]ghost.Transaction, 0), 1}
	// Validate created block
	for i := 0; ; i++ {
		genesisBlock.Nonce = i
		if !blockchain.IsHashValid(blockchain.CalculateHash(genesisBlock), genesisBlock.Difficulty) {
			continue
		} else {
			genesisBlock.Hash = blockchain.CalculateHash(genesisBlock)
			break
		}

	}
	mutex.Lock()
	blockchain.CurrentBlockchain.Blocks = append(blockchain.CurrentBlockchain.Blocks, genesisBlock)
	mutex.Unlock()

	// Create a first node in the network
	alice, err := noise.NewNode()
	check(err)

	bob, err := noise.NewNode()
	check(err)

	// TODO: Adding the way to handle the incoming blocks for the blockchain
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
