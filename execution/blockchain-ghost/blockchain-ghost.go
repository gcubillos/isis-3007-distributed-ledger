package main

import (
	"fmt"
	ghost "github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain-ghost"
	"strconv"
)

// *** Execution of small scale tests ***
func main() {
	// Creating an array of possible users
	var users []string
	users = make([]string, 0, 0)
	for i := 0; i < 10; i++ {
		users = append(users, strconv.Itoa(i))
	}
	// Creating accounts for the users
	for i := range users {
		ghost.CreateAccount("123" + strconv.Itoa(i))
	}
	fmt.Printf("", users)
	// Creating new nodes
	ghost.GenerateNode()

	//// Creating a network node
	//nodeA := ghost.GenerateNode()
	//
	//// Creating genesis state
	//
	//// Creating network with no blocks and capacity 1
	//
	///* Creating the genesis block with the starting parameters for the network
	// */
	//var theBlock = generateBlock(time.Now(), 1, nil, nil, nil)
	//
	//testGhost.blocks = make([]block, 0, 1)
	//
	//// Creating test accounts
	//var testAccount1 = account{
	//	nonce:   0,
	//	balance: 0,
	//	address: "172",
	//}
	//var testAccount2 = account{
	//	nonce:   0,
	//	balance: 4,
	//	address: "174",
	//}
	//
	//// Creating a state
	//var testState = make(map[string]*account)
	//// Inputting values
	//testState["0"] = &testAccount1
	//testState["1"] = &testAccount2
	//
	//// Creating a transaction
	//
	//// Testing state transition

}
