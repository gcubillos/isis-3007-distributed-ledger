package blockchain

import (
	"github.com/perlin-network/noise"
	"sync"
)

// Mutual exclusion variable
var mutex = &sync.Mutex{}

// What the node contains, the data structure and a reference to a peer in the p2p network
type BlockchainNode struct {
	DataStructure Blockchain
	Node          *noise.Node
}

// Create a node in the network
// The genesis block is passed to the Node and a first peer to connect to the network
//func CreateNode (pGenesisBlock Block, pNode noise.Node) BlockchainNode {
//
//	return
//}
