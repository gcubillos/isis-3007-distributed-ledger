package main

import (
	"fmt"
	"github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain"
	"github.com/gcubillos/isis-3007-distributed-ledger/data-structures/shared-components"
	"time"
)

func main() {

	// Defining parameters for simple execution

	// Defining the difficulty for the tests (number of leading zeroes required in the hash)
	var definedDifficulty = 1

	// Defining the amount of currency that will be available during the tests.
	// Of course, transactions with 0 value can be made as well. This is for the sake of simplicity, to have
	// a fixed amount of currency
	var availableCurrency = 10.0
	t := time.Now()
	// TODO: Constructor for genesis blocks
	genesisBlock := blockchain.Block{Timestamp: t, Hash: "", Transactions: make([]components.Transaction, 0), Difficulty: definedDifficulty}
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
	otherNode := blockchain.CreateNode(firstNode.DataStructure, firstNode.Node)

	// Create an example transaction
	exampleTransaction := components.CreateTransaction(firstNode.Node.Addr(), firstNode.Node.Addr(), otherNode.Node.Addr(), 0)

	// Include transaction in list
	transactionList := make([]components.Transaction, 1, 1)
	transactionList[0] = exampleTransaction

	// Throughput tests
	startingBlock := genesisBlock

	// Count number of blocks
	i := 0
	for i = 0; i < 1; i++{
		startingTime := time.Now()
		startingBlock = otherNode.GenerateBlock(startingBlock, transactionList)
		finishTime := time.Now()
		times := finishTime.Sub(startingTime)
		fmt.Printf("%v \n starting %v \n finishing %v",times, startingTime,finishTime)

	}
	fmt.Printf("The number of blocks generated were: %v", i)

}
