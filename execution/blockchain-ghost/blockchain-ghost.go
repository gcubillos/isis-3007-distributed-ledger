package main

import (
	"fmt"
	ghost "github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain-ghost"
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
	genesisBlock := ghost.Block{
		Timestamp:         t,
		Hash:              "",
		HashPreviousBlock: "",
		Parent:            nil,
		Uncles:            nil,
		Transactions:      make([]components.Transaction, 0),
		RecentState:       make(map[string]*ghost.Account, 0),
		BlockNumber:       0,
		Difficulty:        definedDifficulty,
	}

	// For simplicity a "main" account will be created that contains the amount of currency available
	mainAccount := ghost.CreateAccount("main")
	mainAccount.Balance = availableCurrency
	genesisBlock.RecentState[mainAccount.Address] = &mainAccount

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
	firstNode := ghost.CreateInitialNode(genesisBlock)

	// Create other nodes
	otherNode := ghost.GenerateNode(firstNode.DataStructure, firstNode.Node)

	// Create an empty transaction
	exampleTransaction := components.CreateTransaction(firstNode.Node.Addr(), firstNode.Node.Addr(), otherNode.Node.Addr(), 0)

	transactionList := make([]components.Transaction, 1, 1)
	transactionList[0] = exampleTransaction

	theFirstBlock := firstNode.GenerateBlock(&genesisBlock, transactionList)

	fmt.Printf("\n %v", theFirstBlock.Hash)

	firstNode.GenerateBlock(&theFirstBlock, transactionList)

	// Latency: Time it takes for the transaction to be accepted by the other nodes
	//One node submits transactions to the network and, for each transaction, the receiving
	// node(s) subtract the initiation timestamp from the completion timestamp. The median and standard deviation are reported.
	//	For the block-lattice:
	// The initiation timestamp is taken when the send block (S) is created.
	// The completion timestamp is taken after the receive block (R) is confirmed.
	//
	// The latency may vary depending on the size of the data structure. For this reason, it is measured for different transactions as each data structure continues to grow.

	// Throughput: To measure throughput: one node submits as many transactions as it
	// can for one second.
	// Each throughput experiment is repeated five times and the median and standard
	// deviation for the five runs are reported.

}
