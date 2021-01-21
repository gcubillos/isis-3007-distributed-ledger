package main

import (
	"github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain"
	ghost "github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain-ghost"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

// The following code tries to perform the tests regarding the metric of latency on the blockchain data structure.

func main() {
	// Defining the difficulty for the tests (number of leading zeroes required in the hash)
	var definedDifficulty = 1
	blockchain.Difficulty = definedDifficulty
	// Defining the amount of currency that will be available during the tests.
	// Of course, transactions with 0 value can be made as well. This is for the sake of simplicity, to have
	// a fixed amount of currency
	var availableCurrency = 10.0
	// Creating the genesis block
	t := time.Now()
	genesisBlock := blockchain.Block{}
	genesisBlock = blockchain.Block{Timestamp: t, Hash: blockchain.CalculateHash(genesisBlock), Transactions: make([]ghost.Transaction, 0), Difficulty: definedDifficulty}
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

	// Create the first node in the network to have as a starting point
	firstNode := blockchain.CreateInitialNode(genesisBlock, availableCurrency)

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
