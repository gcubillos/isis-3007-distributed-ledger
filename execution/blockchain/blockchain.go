package main

import (
	"github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain"
	ghost "github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain-ghost"
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

	// Create the first node in the network to have as a starting point
	firstNode := blockchain.CreateInitialNode(genesisBlock)

	// Create other nodes
	otherNode := blockchain.CreateNode(genesisBlock, firstNode.Node)

	// Create an empty transaction
	exampleTransaction := ghost.Transaction{
		Origin:          firstNode.Node.Addr(),
		SenderSignature: firstNode.Node.Addr(),
		Destination:     otherNode.Node.Addr(),
		Value:           0,
	}
	transactionList := make([]ghost.Transaction, 1, 1)
	transactionList[0] = exampleTransaction

	firstNode.GenerateBlock(genesisBlock, transactionList)

}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
