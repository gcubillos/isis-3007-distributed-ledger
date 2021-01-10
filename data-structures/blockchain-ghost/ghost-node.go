package ghost

import (
	"github.com/perlin-network/noise"
	"time"
)

// Declaration of node in the network
// Contains the underlying data structure as well as the node from the noise library
type ghostNode struct {
	dataStructure ghost
	node          *noise.Node
}

// *** Constructors ***

// Creating a node in the network
func GenerateNode() ghostNode {
	var rNode = ghostNode{}
	node, err := noise.NewNode()
	check(err)
	rNode.node = node
	rNode.dataStructure = ghost{blocks: []block{}, state: make(map[string]*account)}
	rNode.mining()
	return rNode
}

// *** Methods ***

/* Creating a standard block in the network
 */
func (*ghostNode) generateBlock(pNonce int, pParent *block,
	pTransactions []transaction, pEndState map[string]*account) block {
	var rBlock block
	rBlock.parent = pParent
	rBlock.timestamp = time.Now()
	rBlock.nonce = pNonce
	rBlock.hashPreviousBlock = calculateHash(*pParent)
	rBlock.transactions = pTransactions
	rBlock.endState = pEndState
	return rBlock
}

// Revises whether the error is not nil
func check(err error) {
	if err != nil {
		panic(err)
	}
}
