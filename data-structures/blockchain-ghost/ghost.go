package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

// *** Structs ***

/* Declaration of structure
Containing the blocks and the genesis state
*/
// TODO: Including state in ghost struct?
type ghost struct {
	blocks []block
	state  map[string]*account
}

// What a block in the network contains
// The network is intended to produce roughly one block every ten minutes, with each block 
// containing a timestamp, a nonce, a reference to (ie. hash of) the previous block and a
// list of all of the transactions that have taken place since the previous block.

type block struct {
	timestamp         time.Time
	nonce             int
	hashPreviousBlock string
	parent            *block
	uncles            []block
	transactions      []transaction
	endState          map[string]*account
}

// What a transaction ensues
// A transaction is a request to move $X from A to B
type transaction struct {
	origin          string
	senderSignature string
	destination     string
	value           float32
}

// What an account contains
// Nonce counter used to make sure each transaction can only be processed once
// account's current balance
type account struct {
	nonce   int
	balance float32
	address string
}

// *** Functions ***

/* Creating a standard block in the network
 */
func generateBlock(pTimestamp time.Time, pNonce int, pParent *block,
	pTransactions []transaction, pEndState map[string]*account) block {
	var rBlock block
	rBlock.parent = pParent
	rBlock.timestamp = pTimestamp
	rBlock.nonce = pNonce
	rBlock.hashPreviousBlock = calculateHash(*pParent)
	rBlock.transactions = pTransactions
	rBlock.endState = pEndState
	return rBlock
}

// Check that the block is valid
// By checking if the previous block referenced by the block exists and is valid
// Checking that the timestamp of the block is greater than that of the previous block
// Check that the proof of work on the block is valid.
// Let S[0] be the state at the end of the previous block.
// Suppose TX is the block's transaction list with n transactions. For all i in 0...n-1,
// set S[i+1] = APPLY(S[i],TX[i]) If any application returns an error, exit and return false.
// Return true, and register S[n] as the state at the end of this block.

func checkBlockValid(pBlock block) (isValid bool) {
	isValid = true
	// Previous block exists and valid
	if !checkBlockValid(*pBlock.parent) {
		isValid = false
	}
	// Timestamp
	if pBlock.timestamp.Before(pBlock.parent.timestamp) {
		isValid = false
	}
	// Proof of work
	// Simplified version of proof of work
	// TODO: Checking proof of work

	// State transition check
	var initialState = pBlock.parent.endState
	for i := 0; i < len(pBlock.transactions) && isValid; i++ {
		currentState, err := stateTransition(initialState, pBlock.transactions[i])
		if err != "" {
			isValid = false
		}
		initialState = currentState
	}
	// Checking state
	if len(pBlock.endState) == len(initialState) {
		for i, _ := range pBlock.endState {
			if !(pBlock.endState[i].nonce == initialState[i].nonce) &&
							(pBlock.endState[i].balance == initialState[i].balance) &&
							(pBlock.endState[i].address == initialState[i].address) {
				isValid = false
				break
			}
		}
	} else {
		isValid = false
	}

	return isValid
}

/* Checks validity of uncles
 */
// TODO: Finish validity of uncles
func checkUncleValidity(pBlock block) (isValid bool) {
	return false
}

// State transition function. Checks validity of a change in state from a list of transactions
// Syntax APPLY(S,TX) -> S'
func stateTransition(pCurrentState map[string]*account, pTransaction transaction) (pModifiedState map[string]*account, err string) {
	// If referenced UTXO is not in S
	err = ""
	pModifiedState = pCurrentState
	if pCurrentState[pTransaction.senderSignature].balance <= pTransaction.value {
		err = "The referenced UTXO is not in the state\n"
	}
	// If the provided signature does not match the owner of the UTXO
	if pTransaction.origin != pTransaction.senderSignature {
		err = err + "The provided signature does not match the owner of the UTXO\n"
	}
	// If the sum of the denominations of all input UTXO is less than the sum of the
	// denominations of all output UTXO, return an error. Not necessary given that a
	// transaction struct only contains one transaction.

	// Return S'. Apply the changes in the transaction
	if err == "" {
		pModifiedState[pTransaction.origin].balance -= pTransaction.value
		pModifiedState[pTransaction.destination].balance += pTransaction.value
	}
	// Creating the account in the state if it doesn't already exist
	if pTransaction.destination not in 
	return pModifiedState, err
}

// Generate hash of a block. Using block header which includes timestamp, nonce,
// previous block hash
func calculateHash(pBlock block) (rHash string) {
	bHeader := strconv.Itoa(pBlock.nonce) + pBlock.timestamp.String() + pBlock.hashPreviousBlock
	hash := sha256.New()
	hash.Write([]byte(bHeader))
	rHash = hex.EncodeToString(hash.Sum(nil))
	return rHash
}

// *** Execution of small scale tests ***
func main() {
	// Creating genesis state


	// Creating network with no blocks and capacity 1
	var testGhost = ghost{make([]block, 0, 1)}

	/* Creating the genesis block with the starting parameters for the network
	 */
	var theBlock = generateBlock(time.Now(),1,nil,nil,nil)

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
	var testState = make(map[string]*account);
	// Inputting values
	testState["0"] = &testAccount1;
	testState["1"] = &testAccount2;


	// Creating a transaction


	// Testing state transition

}
