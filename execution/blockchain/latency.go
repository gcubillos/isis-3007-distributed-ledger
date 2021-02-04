package main

import (
	"github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain"
	"github.com/gcubillos/isis-3007-distributed-ledger/data-structures/shared-components"
	"math/rand"
	"time"
)

// The following code tries to perform the tests regarding the metric of latency on the blockchain data structure.

func main() {
	// TODO: Making version that asks the users through console what parameters they want
	// TODO: Using defer for closing the nodes?
	// Defining parameters for testing

	// Defining number of nodes to be present in the network additionally to the initial node
	var numberNodes = 10

	// Defining number of transactions to occur in the network
	var numberTransactions = 10

	// Defining the difficulty for the tests (number of leading zeroes required in the hash)
	var definedDifficulty = 10

	// Defining the amount of currency that will be available during the tests.
	// Of course, transactions with 0 value can be made as well. This is for the sake of simplicity, to have
	// a fixed amount of currency
	var availableCurrency = 10.0

	// Creating the genesis block
	t := time.Now()
	genesisBlock := blockchain.Block{}
	genesisBlock = blockchain.Block{Timestamp: t, Hash: blockchain.CalculateHash(genesisBlock), Transactions: make([]components.Transaction, 0), Difficulty: definedDifficulty}
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
	// Assign difficulty
	blockchain.Difficulty = definedDifficulty

	// Create the first node in the network to have as a starting point
	firstNode := blockchain.CreateInitialNode(genesisBlock, availableCurrency)

	// Array for keeping track of the nodes' addresses without having to ask the network
	var nodesNetwork []blockchain.NodeBlockchain
	nodesNetwork = make([]blockchain.NodeBlockchain, 0)

	// Create other nodes
	for i := 0; i < numberNodes; i++ {
		nodesNetwork = append(nodesNetwork, blockchain.CreateNode(firstNode.DataStructure, firstNode.Node))
	}

	// Creating seed for randomizing sender and receiver of transactions
	rand.Seed(time.Now().UnixNano())

	// Create transactions
	for j := 0; j < numberTransactions; j++ {
		var randomSender int
		var randomReceiver int
		for {
			randomSender = rand.Intn(len(nodesNetwork))
			randomReceiver = rand.Intn(len(nodesNetwork))
			if randomReceiver != randomSender {
				break
			}
		}

		exampleTransaction := components.Transaction{
			Origin:          nodesNetwork[randomSender].Node.Addr(),
			SenderSignature: nodesNetwork[randomSender].Node.Addr(),
			Destination:     nodesNetwork[randomReceiver].Node.Addr(),
			Value:           rand.Float64(),
		}
		transactionList := make([]components.Transaction, 1, 1)
		transactionList[0] = exampleTransaction

		nodesNetwork[rand.Intn(len(nodesNetwork))].GenerateBlock(genesisBlock, transactionList)
	}

	// TODO: Make the tests make a little bit more sense, such that each node runs them individually, not a local machine

}
