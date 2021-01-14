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
	theBlockchain := blockchain.Blockchain{
		Blocks: []blockchain.Block{},
		State:  make(map[string]float64, 0),
	}
	mutex.Lock()
	theBlockchain.Blocks = append(theBlockchain.Blocks, genesisBlock)
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

		anotherBlockchain := blockchain.Blockchain{
			Blocks: make([]blockchain.Block, 0),
			State:  make(map[string]float64, 0),
		}
		if err := json.Unmarshal(ctx.Data(), &anotherBlockchain); err != nil {
			fmt.Printf("", err.Error())
		}
		fmt.Printf("", anotherBlockchain.Blocks)

		return nil
	})

	alice.Handle(func(ctx noise.HandlerContext) error {
		if !ctx.IsRequest() {
			return nil
		}

		return nil
	})

	err = alice.Listen()
	check(err)
	err = bob.Listen()
	check(err)
	// Creating a transaction
	exampleTransaction := ghost.Transaction{
		Origin:          alice.Addr(),
		SenderSignature: alice.ID().Address,
		Destination:     bob.Addr(),
		Value:           5,
	}
	transactionList := make([]ghost.Transaction, 1, 1)
	transactionList[0] = exampleTransaction
	blockchain.CurrentBlockchain = theBlockchain
	blockchain.CurrentBlockchain.State[alice.Addr()] = 5
	theBlock := blockchain.GenerateBlock(genesisBlock, transactionList)
	if v, err := blockchain.IsBlockValid(theBlock, genesisBlock); v {
		mutex.Lock()
		blockchain.CurrentBlockchain.Blocks = append(blockchain.CurrentBlockchain.Blocks, theBlock)
		mutex.Unlock()
	} else {
		panic(err)
	}

	bytes, err := json.Marshal(blockchain.CurrentBlockchain)
	res, err := alice.Request(context.TODO(), bob.Addr(), bytes)

	fmt.Printf("", res)

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
