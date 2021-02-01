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
var theNode NodeGhost

// *** Constructors ***

// Create a node in the network such that it can discover other nodes using the Kademlia
// protocol. The current state of the blockchain is passed to the Node and a first peer
// to connect to the network
func GenerateNode(pCurrentGhost Ghost, pNode *noise.Node) NodeGhost {
	// Create structure
	theNode = NodeGhost{
		DataStructure: pCurrentGhost,
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

		receivedBlockchain := Ghost{
			Blocks: make([]Block, 0),
			State:  make(map[string]*Account, 0),
		}
		if err := json.Unmarshal(ctx.Data(), &receivedBlockchain); err != nil {
			fmt.Printf("trouble unmarshalling CreateNode. Error: %v Blockchain: %v \n", err.Error(), receivedBlockchain.Blocks)
		} else {
			theNode.DataStructure.ReplaceGHOST(receivedBlockchain)
		}
		fmt.Printf("current structure CreateNode %v \n", theNode.DataStructure)

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
	theNode.Node = networkNode

	return theNode
}

// *** Methods ***

/* Creating a standard Block in the network
 */
func (*NodeGhost) generateBlock(pNonce int, pParent *Block,
	pTransactions []components.Transaction, pEndState map[string]*Account) Block {
	var rBlock Block
	rBlock.Parent = pParent
	rBlock.Timestamp = time.Now()
	rBlock.Nonce = pNonce
	rBlock.HashPreviousBlock = pParent.calculateHash()
	rBlock.Transactions = pTransactions
	rBlock.RecentState = pEndState
	// Proof of work
	// TODO: Including Simplified version of proof of work
	return rBlock
}

// Revises whether the error is not nil
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Discover other nodes in the network
func (a *NodeGhost) discover() {
}
