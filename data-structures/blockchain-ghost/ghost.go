package main

import "time"

// *** Structs ***

/* Declaration of structure
Containing the blocks and the initial state
 */
type ghost struct {
	blocks []block
	state map[int]account
}

// What a block in the network contains
// The network is intended to produce roughly one block every ten minutes, with each block 
// containing a timestamp, a nonce, a reference to (ie. hash of) the previous block and a
// list of all of the transactions that have taken place since the previous block.

type block struct {
	timestamp time.Time
	nonce int
	hashPreviousBlock string
	parent *block
	uncles []block
	transactions []transaction
}

// What a transaction ensues
// A transaction is a request to move $X from A to B
type transaction struct {
	origin string
	destination string
	amount float32
}

// What an account contains
type account struct{
	nonce int
	balance float32
}

// *** Functions ***

/* Creating a standard block in the network
 */
func generateBlock(pTimestamp time.Time, pNonce int, pParent block) (rBlock block) {

	return rBlock
}

// Check that the block is valid
// By checking if the previous block referenced by the block exists and is valid
// Checking that the timestamp of the block is greater than that of the previous block
// Check that the proof of work on the block is valid.
// Let S[0] be the state at the end of the previous block.
// Suppose TX is the block's transaction list with n transactions. For all i in 0...n-1, set S[i+1] = APPLY(S[i],TX[i]) If any application returns an error, exit and return false.
// Return true, and register S[n] as the state at the end of this block.

func checkBlockValid(pBlock block) (isValid bool) {
	isValid = false
	if pBlock.nonce < 1 {
		isValid = true
	}
	return
}

/* Checks validity of uncles
 */
func checkUncleValidity(pBlock block) (isValid bool) {
	return false
}

// State transition function. Checks validity of a change in state from a list of transactions
// Syntax APPLY(S,TX) -> S'
func stateTransition (pState map[string]int, pTransaction transaction) (isValid bool){
	// If referenced UTXO  is not in S
	if {}
}

// *** Execution of small scale tests ***
func main(){
	// Creating network with no blocks and capacity 1
	var testGhost = ghost{make([]block,0,1),make(map[int]account)}
	/* Creating the genesis block with the starting parameters for the network
	 */
	var genesisBlock =
	var theGhost = new(ghost)
	var theBlock = new(block)
	theBlock.nonce = 2
	theGhost.blocks = make([]block,0,1)
	theGhost.blocks = append(theGhost.blocks, *theBlock)

}
