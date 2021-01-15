package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	ghost "github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain-ghost"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
	"sync"
	"time"
)

// Mutual exclusion variable
var mutex = &sync.Mutex{}

// An instance of the node blockchain
var theNode NodeBlockchain

// What the node contains, the data structure and a reference to a peer in the p2p network
type NodeBlockchain struct {
	DataStructure Blockchain
	Node          *noise.Node
}

// Create a node in the network such that it can discover other nodes using the Kademlia protocol
// The genesis block is passed to the Node and a first peer to connect to the network
func CreateNode(pGenesisBlock Block, pNode *noise.Node) NodeBlockchain {
	// Create structure
	theNode = NodeBlockchain{
		DataStructure: Blockchain{[]Block{pGenesisBlock}, make(map[string]float64, 0)},
		Node:          nil,
	}
	// Create network node
	networkNode, err := noise.NewNode()
	check(err)

	// Assign the Kademlia protocol to the node so it can discover other nodes
	ka := kademlia.New()
	networkNode.Bind(ka.Protocol())

	// Assign the way the node will handle the requests for blockchain updates
	networkNode.Handle(func(ctx noise.HandlerContext) error {
		if !ctx.IsRequest() {
			return nil
		}

		receivedBlockchain := Blockchain{
			Blocks: make([]Block, 0),
			State:  make(map[string]float64, 0),
		}
		if err := json.Unmarshal(ctx.Data(), &receivedBlockchain); err != nil {
			fmt.Printf("trouble unmarshalling", err.Error())
		} else {
			theNode.DataStructure.ReplaceChain(receivedBlockchain)
		}
		fmt.Printf("current structure", theNode.DataStructure)

		return nil
	})

	// Make the node listen to the network
	check(networkNode.Listen())

	// Ping the provided node in the network
	networkNode.Ping(context.TODO(), pNode.Addr())
	// Discover the other nodes present in the network at the moment
	ka.Discover()
	// Assign the network node to the node
	theNode.Node = networkNode

	return theNode
}

// Create the initial node
// The genesis block is passed to the Node
func CreateInitialNode(pGenesisBlock Block) NodeBlockchain {
	// Create structure
	theNode = NodeBlockchain{
		DataStructure: Blockchain{[]Block{pGenesisBlock}, make(map[string]float64, 0)},
		Node:          nil,
	}
	// Create network node
	networkNode, err := noise.NewNode()
	check(err)

	// Assign the Kademlia protocol to the node so it can discover other nodes
	ka := kademlia.New()
	networkNode.Bind(ka.Protocol())

	// Assign the way the node will handle the requests for blockchain updates
	networkNode.Handle(func(ctx noise.HandlerContext) error {
		if !ctx.IsRequest() {
			return nil
		}

		receivedBlockchain := Blockchain{
			Blocks: make([]Block, 0),
			State:  make(map[string]float64, 0),
		}
		if err := json.Unmarshal(ctx.Data(), &receivedBlockchain); err != nil {
			fmt.Printf("trouble unmarshalling", err.Error())
		} else {
			theNode.DataStructure.ReplaceChain(receivedBlockchain)
		}
		fmt.Printf("current structure", theNode.DataStructure)

		return nil
	})

	// Make the node listen to the network
	check(networkNode.Listen())

	// Assign the network node to the node
	theNode.Node = networkNode

	return theNode
}

// Create a block and broadcast it to the rest of the network

func (pNode *NodeBlockchain) GenerateBlock(oldBlock Block, pTransactions []ghost.Transaction) Block {

	var newBlock Block

	newBlock.Timestamp = time.Now()
	newBlock.Transactions = pTransactions
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty
	// Calculating the hash
	for i := 0; ; i++ {
		newBlock.Nonce = i
		if !IsHashValid(CalculateHash(newBlock), newBlock.Difficulty) {
			continue
		} else {
			newBlock.Hash = CalculateHash(newBlock)
			break
		}
	}

	// Check that the block is valid
	if ok, err := theNode.DataStructure.IsBlockValid(newBlock, oldBlock); ok {
		check(err)
		// Add the block to the current blockchain
		mutex.Lock()
		theNode.DataStructure.Blocks = append(theNode.DataStructure.Blocks, newBlock)
		mutex.Unlock()
		// Convert the blockchain so that it can be sent
		bytes, err := json.Marshal(theNode.DataStructure)
		check(err)
		// Broadcast the blockchain to the network
		for _, v := range theNode.Node.Outbound() {
			theNode.Node.Send(context.TODO(), v.ID().Address, bytes)
		}
	}

	return newBlock
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
