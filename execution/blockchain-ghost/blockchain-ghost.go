package main

import (
	ghost "github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain-ghost"
	"github.com/gcubillos/isis-3007-distributed-ledger/data-structures/shared-components"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

func main() {

	// Defining parameters for simple execution

	// Defining the difficulty for the tests (number of leading zeroes required in the hash)
	var definedDifficulty = 1

	// Defining the amount of currency that will be available during the tests.
	// Of course, transactions with 0 value can be made as well. This is for the sake of simplicity, to have
	// a fixed amount of currency
	var availableCurrency = 10.0
	t := time.Now()
	genesisBlock := ghost.Block{}
	genesisBlock = ghost.Block{
		Timestamp:         t,
		Hash:              ghost.CalculateHash(genesisBlock),
		HashPreviousBlock: "",
		Parent:            nil,
		Uncles:            nil,
		Transactions:      make([]components.Transaction, 0),
		RecentState:       nil,
		BlockNumber:       0,
		Difficulty:        0,
	}
	genesisBlock = ghost.Block{Timestamp: t, Hash: ghost.CalculateHash(genesisBlock), Transactions: make([]components.Transaction, 0), Difficulty: definedDifficulty}
	// Validate created block
	for i := 0; ; i++ {
		genesisBlock.Nonce = i
		if !ghost.IsHashValid(ghost.CalculateHash(genesisBlock), genesisBlock.Difficulty) {
			continue
		} else {
			genesisBlock.Hash = ghost.CalculateHash(genesisBlock)
			break
		}

	}

	// Create the first node in the network to have as a starting point
	firstNode := ghost.CreateInitialNode(genesisBlock, availableCurrency)

	// Create other nodes
	otherNode := ghost.CreateNode(firstNode.DataStructure, firstNode.Node)

	// Create an empty transaction
	exampleTransaction := components.Transaction{
		Origin:          firstNode.Node.Addr(),
		SenderSignature: firstNode.Node.Addr(),
		Destination:     otherNode.Node.Addr(),
		Value:           0,
	}
	transactionList := make([]components.Transaction, 1, 1)
	transactionList[0] = exampleTransaction

	firstNode.GenerateBlock(genesisBlock, transactionList)

}
