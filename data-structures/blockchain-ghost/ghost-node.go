package ghost

import (
	"github.com/gcubillos/isis-3007-distributed-ledger/data-structures/shared-components"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
	"time"
)

// Declaration of node in the network
// Contains the underlying data structure as well as the node from the noise library
type NodeGhost struct {
	DataStructure Ghost
	Node          *noise.Node
}

// *** Constructors ***

// Creating a node in the network
// Bind with the Kademlia protocol so that it can discover other peers
// Also makes the node listen for other connections
func GenerateNode() NodeGhost {
	var rNode = NodeGhost{}
	node, err := noise.NewNode()
	check(err)
	rNode.Node = node
	rNode.DataStructure = Ghost{Blocks: []Block{}, State: make(map[string]*Account)}
	commProtocol := kademlia.New()
	rNode.Node.Bind(commProtocol.Protocol())
	if err := rNode.Node.Listen(); err != nil {
		panic(err)
	}
	return rNode
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
