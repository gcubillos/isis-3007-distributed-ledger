package main

import (
	ghost "github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain-ghost"
	"time"
)

// *** Execution of small scale tests ***
func main() {
	// Creating a network node
	nodeA := ghost.GenerateNode()

	// Creating genesis state

	// Creating network with no blocks and capacity 1

	/* Creating the genesis block with the starting parameters for the network
	 */
	var theBlock = generateBlock(time.Now(), 1, nil, nil, nil)

	testGhost.blocks = make([]block, 0, 1)

	// Creating test accounts
	var testAccount1 = account{
		nonce:   0,
		balance: 0,
		address: "172",
	}
	var testAccount2 = account{
		nonce:   0,
		balance: 4,
		address: "174",
	}

	// Creating a state
	var testState = make(map[string]*account)
	// Inputting values
	testState["0"] = &testAccount1
	testState["1"] = &testAccount2

	// Creating a transaction

	// Testing state transition

}
