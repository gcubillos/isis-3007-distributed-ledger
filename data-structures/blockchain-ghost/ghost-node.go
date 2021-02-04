package ghost

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gcubillos/isis-3007-distributed-ledger/data-structures/shared-components"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
	"sync"
	"time"
)

// Mutual exclusion variable
var mutex = &sync.Mutex{}

// Declaration of node in the network
// Contains the underlying data structure as well as the node from the noise library
type NodeGhost struct {
	DataStructure Ghost
	Node          *noise.Node
}

// An instance of the ghost node
var thisNode NodeGhost

// *** Constructors ***

// Create a node in the network such that it can discover other nodes using the Kademlia
// protocol. The current state of the blockchain is passed to the Node and a first peer
// to connect to the network
func GenerateNode(pCurrentGhost Ghost, pNode *noise.Node) NodeGhost {
	// Create structure
	thisNode = NodeGhost{
		DataStructure: pCurrentGhost,
		Node:          nil,
	}
	// Create network node
	networkNode, err := noise.NewNode()
	check(err)

	// Assign the Kademlia protocol to the node so it can discover other nodes
	ka := kademlia.New()
	networkNode.Bind(ka.Protocol())

	// Assign the way the node will handle the requests for updates in the chain
	networkNode.Handle(func(ctx noise.HandlerContext) error {
		if !ctx.IsRequest() {
			return nil
		}

		receivedGhost := Ghost{
			Blocks: make([]Block, 0),
			State:  make(map[string]*Account, 0),
		}
		if err := json.Unmarshal(ctx.Data(), &receivedGhost); err != nil {
			fmt.Printf("trouble unmarshalling CreateNode. Error: %v Blockchain: %v \n", err.Error(), receivedGhost.Blocks)
		} else {
			thisNode.DataStructure.FindGHOST(receivedGhost)
		}
		fmt.Printf("current structure CreateNode %v \n", thisNode.DataStructure)

		return nil
	})

	// Make the node listen to the network
	check(networkNode.Listen())

	// Ping the provided node in the network
	_, err = networkNode.Ping(context.TODO(), pNode.Addr())
	check(err)

	// Discover the other nodes present in the network at the moment
	ka.Discover()
	// Assign the network node to the node
	thisNode.Node = networkNode

	return thisNode
}

// Create the initial node
// The genesis block is passed to the Node
// The amount of available currency is passed to the node
func CreateInitialNode(pGenesisBlock Block, pAvailableCurrency float64) NodeGhost {
	// Create structure
	thisNode = NodeGhost{
		DataStructure: Ghost{[]Block{pGenesisBlock}, make(map[string]*Account, 0), []Block{pGenesisBlock}},
		Node:          nil,
	}
	// Create network node
	networkNode, err := noise.NewNode()
	check(err)

	// For simplicity a "main" account will be created that contains the amount of currency available
	mainAccount := CreateAccount(networkNode.Addr())
	mainAccount.Balance = pAvailableCurrency
	thisNode.DataStructure.State[networkNode.Addr()] = &mainAccount

	// Assign the Kademlia protocol to the node so it can discover other nodes
	ka := kademlia.New()
	networkNode.Bind(ka.Protocol())

	// Assign the way the node will handle the requests for updates in the chain
	networkNode.Handle(func(ctx noise.HandlerContext) error {
		if !ctx.IsRequest() {
			return nil
		}

		receivedGhost := Ghost{
			Blocks: make([]Block, 0),
			State:  make(map[string]*Account, 0),
		}
		// TODO: Avoid having the unmarshal error when discovering peers. Check the kademlia discover method
		if err := json.Unmarshal(ctx.Data(), &receivedGhost); err != nil {
			fmt.Printf("trouble unmarshalling InitialNode. Error: %v Blockchain: %v \n", err.Error(), receivedGhost)
		} else {
			thisNode.DataStructure.FindGHOST(receivedGhost)
		}
		fmt.Printf("current structure InitialNode %v \n", thisNode.DataStructure)

		return nil
	})

	// Make the node listen to the network
	check(networkNode.Listen())

	// TODO: Modify the time the address is being put to the state so that the nodes can communicate

	// Assign the network node to the node
	thisNode.Node = networkNode

	return thisNode
}

// *** Methods ***

// Creating a standard Block in the network and broadcasting it
func (pNode *NodeGhost) GenerateBlock(pParent *Block, pTransactions []components.Transaction) Block {

	var nBlock Block

	// Adding the transaction that gives the "miner" a reward for doing the work
	// TODO: Revise the rewards whether it is belonging to the main chain or not
	rewardTransaction := components.Transaction{
		Origin:          "main",
		SenderSignature: "main",
		Destination:     pNode.Node.Addr(),
		Value:           1,
	}
	pTransactions = append(pTransactions, rewardTransaction)

	// Basic information in the block
	nBlock.Parent = pParent
	nBlock.Timestamp = time.Now()
	nBlock.HashPreviousBlock = pParent.Hash
	nBlock.Difficulty = pParent.Difficulty
	nBlock.Transactions = pTransactions
	nBlock.BlockNumber = len(pNode.DataStructure.Blocks) + 1

	// Proof of work, calculating the hash
	for i := 0; ; i++ {
		nBlock.Nonce = i
		if !IsHashValid(CalculateHash(nBlock), nBlock.Difficulty) {
			continue
		} else {
			nBlock.Hash = CalculateHash(nBlock)
			break
		}
	}

	// Check that the block is valid
	if ok, err := thisNode.DataStructure.IsBlockValid(nBlock); ok {
		check(err)
		// Add the block to the current structure
		mutex.Lock()
		thisNode.DataStructure.Blocks = append(thisNode.DataStructure.Blocks, nBlock)
		thisNode.DataStructure.CurrentChain = append(thisNode.DataStructure.CurrentChain, nBlock)
		mutex.Unlock()
		// Convert the chain so that it can be broadcast
		bytes, err := json.Marshal(thisNode.DataStructure)
		check(err)
		// Broadcast the chain to the network
		for _, v := range thisNode.Node.Outbound() {
			_, err = thisNode.Node.Request(context.TODO(), v.ID().Address, bytes)
			check(err)
		}
	} else {

	}

	return nBlock
}

// Revises whether the error is not nil
func check(err error) {
	if err != nil {
		panic(err)
	}
}
