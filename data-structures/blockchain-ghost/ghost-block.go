package ghost

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

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

// *** Constructors ***

// *** Methods ***
// Check that the block is valid
// By checking if the previous block referenced by the block exists and is valid
// Checking that the timestamp of the block is greater than that of the previous block
// Check that the proof of work on the block is valid.
// Let S[0] be the state at the end of the previous block.
// Suppose TX is the block's transaction list with n transactions. For all i in 0...n-1,
// set S[i+1] = APPLY(S[i],TX[i]) If any application returns an error, exit and return false.
// Return true, and register S[n] as the state at the end of this block.

func (pBlock *block) checkBlockValid() (isValid bool) {
	isValid = true
	// Previous block exists and valid
	if !pBlock.parent.checkBlockValid() {
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

//Checks validity of uncles
//TODO: Finish validity of uncles
func (pBlock *block) checkUncleValidity() (isValid bool) {
	return false
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
