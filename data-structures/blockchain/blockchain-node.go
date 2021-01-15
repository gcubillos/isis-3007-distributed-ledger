package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
	"sync"
)

// Mutual exclusion variable
var mutex = &sync.Mutex{}

// What the node contains, the data structure and a reference to a peer in the p2p network
type BlockchainNode struct {
	DataStructure Blockchain
	Node          *noise.Node
}

// Create a node in the network such that it can discover other nodes using the Kademlia protocol
// The genesis block is passed to the Node and a first peer to connect to the network
func CreateNode(pGenesisBlock Block, pNode noise.Node) BlockchainNode {
	// Create structure
	theNode := BlockchainNode{
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
			fmt.Printf("", err.Error())
		} else {
			theNode.DataStructure.ReplaceChain(receivedBlockchain)
		}

		return nil
	})

	// Make the node listen to the network
	check(networkNode.Listen())

	// Ping the provided node in the network
	networkNode.Ping(context.TODO(), pNode.Addr())
	// Discover the other nodes present in the network at the moment

	return theNode
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
