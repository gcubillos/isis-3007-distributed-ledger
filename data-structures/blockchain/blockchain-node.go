package blockchain

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

// What the node contains, the data structure and a reference to a peer in the p2p network
type NodeBlockchain struct {
	DataStructure Blockchain
	Node          *noise.Node
}

// An instance of the blockchain node
var thisNode NodeBlockchain

// Create a node in the network such that it can discover other nodes using the Kademlia
// protocol. The current state of the blockchain is passed to the Node and a first peer
// to connect to the network
func CreateNode(pCurrentBlockchain Blockchain, pNode *noise.Node) NodeBlockchain {
	// Create structure
	thisNode = NodeBlockchain{
		DataStructure: pCurrentBlockchain,
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
		// TODO: Avoid having the unmarshal error when discovering peers. Check the kademlia discover method.
		// Just change the context received. Uncomment to view the error
		if err := json.Unmarshal(ctx.Data(), &receivedBlockchain); err == nil {
			thisNode.DataStructure.ReplaceChain(receivedBlockchain)
			fmt.Printf("current structure Create Node \n")
			for _, v := range thisNode.DataStructure.Blocks {
				fmt.Printf("a block %v \n", v)
			}
		} else {
			// fmt.Printf("trouble unmarshalling CreateNode. Error: %v Blockchain: %v \n", err.Error(), receivedBlockchain.Blocks)
		}

		return ctx.Send([]byte(""))
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
// The amount of available currency is passed as well to the node
func CreateInitialNode(pGenesisBlock Block, pAvailableCurrency float64) NodeBlockchain {
	// Create structure
	thisNode = NodeBlockchain{
		DataStructure: Blockchain{[]Block{pGenesisBlock}, make(map[string]float64, 0)},
		Node:          nil,
	}
	// For simplicity a "main" account will be created that contains the amount of currency available
	thisNode.DataStructure.State["main"] = pAvailableCurrency
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
		// TODO: Avoid having the unmarshal error when discovering peers. Check the kademlia discover method.
		// Just change the context received. Uncomment to view the error
		if err := json.Unmarshal(ctx.Data(), &receivedBlockchain); err == nil {
			thisNode.DataStructure.ReplaceChain(receivedBlockchain)
			fmt.Printf("current structure InitialNode \n")
			for _, v := range thisNode.DataStructure.Blocks {
				fmt.Printf("a block %v \n", v)
			}
		} else {
			// fmt.Printf("trouble unmarshalling CreateNode. Error: %v Blockchain: %v \n", err.Error(), receivedBlockchain.Blocks)
		}

		return ctx.Send([]byte(""))
	})

	// Make the node listen to the network
	check(networkNode.Listen())

	// Assign the network node to the node
	thisNode.Node = networkNode

	return thisNode
}

// Create a block and broadcast it to the rest of the network

func (pNode *NodeBlockchain) GenerateBlock(oldBlock Block, pTransactions []components.Transaction) Block {

	var newBlock Block

	// Adding the transaction that gives the "miner" a reward for doing the work
	rewardTransaction := components.Transaction{
		Origin:          "main",
		SenderSignature: "main",
		Destination:     pNode.Node.Addr(),
		Value:           1,
	}
	pTransactions = append(pTransactions, rewardTransaction)

	// Including information relevant to the block
	newBlock.Timestamp = time.Now()
	newBlock.Transactions = pTransactions
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = Difficulty
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
	if ok, err := thisNode.DataStructure.IsBlockValid(newBlock, oldBlock); ok {
		check(err)
		// Add the block to the current blockchain
		mutex.Lock()
		thisNode.DataStructure.Blocks = append(thisNode.DataStructure.Blocks, newBlock)
		mutex.Unlock()
		// Convert the blockchain so that it can be sent
		bytes, err := json.Marshal(thisNode.DataStructure)
		check(err)
		// Broadcast the blockchain to the network
		for _, v := range thisNode.Node.Outbound() {
			_, err = thisNode.Node.Request(context.TODO(), v.ID().Address, bytes)
			check(err)
		}
	}

	return newBlock
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
