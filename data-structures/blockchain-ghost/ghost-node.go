package ghost

import (
	"github.com/perlin-network/noise"
	"time"
)

// Declaration of node in the network
// Contains the underlying data structure as well as the node from the noise library
type NodeGhost struct {
	DataStructure ghost
	Node          *noise.Node
}

// *** Constructors ***

// Creating a node in the network
func GenerateNode() NodeGhost {
	var rNode = NodeGhost{}
	node, err := noise.NewNode()
	check(err)
	rNode.Node = node
	rNode.DataStructure = ghost{Blocks: []Block{}, state: make(map[string]*Account)}
	rNode.mining()
	return rNode
}

// *** Methods ***

/* Creating a standard Block in the network
 */
func (*NodeGhost) generateBlock(pNonce int, pParent *Block,
	pTransactions []Transaction, pEndState map[string]*Account) Block {
	var rBlock Block
	rBlock.Parent = pParent
	rBlock.Timestamp = time.Now()
	rBlock.Nonce = pNonce
	rBlock.HashPreviousBlock = pParent.calculateHash()
	rBlock.Transactions = pTransactions
	rBlock.EndState = pEndState
	return rBlock
}

// Revises whether the error is not nil
func check(err error) {
	if err != nil {
		panic(err)
	}
}
