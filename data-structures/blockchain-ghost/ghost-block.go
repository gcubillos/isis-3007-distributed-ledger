package ghost

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"time"
)

// What a Block in the network contains
// The network is intended to produce roughly one Block every ten minutes, with each Block
// containing a Timestamp, a Nonce, a reference to (ie. Hash of) the previous Block and a
// list of all of the Transactions that have taken place since the previous Block.

type Block struct {
	Timestamp         time.Time
	Nonce             int
	HashPreviousBlock string
	Parent            *Block
	Uncles            []Block
	Transactions      []Transaction
	EndState          map[string]*Account
}

// *** Constructors ***

// *** Methods ***
// Check that the Block is valid
// By checking if the previous Block referenced by the Block exists and is valid
// Checking that the Timestamp of the Block is greater than that of the previous Block
// Check that the proof of work on the Block is valid.
// Let S[0] be the state at the end of the previous Block.
// Suppose TX is the Block's Transactions list with n Transactions. For all i in 0...n-1,
// set S[i+1] = APPLY(S[i],TX[i]) If any application returns an error, exit and return false.
// Return true, and register S[n] as the state at the end of this Block.

func (pBlock *Block) IsBlockValid() (bool, error) {
	// Previous Block exists and valid
	if condition, _ := pBlock.Parent.IsBlockValid(); !condition {
		return false, errors.New("previous Block isn't valid")
	}
	// Timestamp
	if pBlock.Timestamp.Before(pBlock.Parent.Timestamp) {
		return false, errors.New("timestamp of previous Block isn't valid")
	}

	// State transition check
	var initialState = pBlock.Parent.EndState
	for i := 0; i < len(pBlock.Transactions); i++ {
		if currentState, err := stateTransition(initialState, pBlock.Transactions[i]); err != nil {
			initialState = currentState
		} else {
			break
		}
	}
	// Checking state
	if len(pBlock.EndState) == len(initialState) {
		for i := range pBlock.EndState {
			if !(pBlock.EndState[i].Nonce == initialState[i].Nonce) &&
				(pBlock.EndState[i].Balance == initialState[i].Balance) &&
				(pBlock.EndState[i].Address == initialState[i].Address) {
				return false, errors.New("there is an error with the state")
			}
		}
	} else {
		return false, errors.New("state doesn't match")
	}

	return true, nil
}

//Checks validity of uncles
//TODO: Finish validity of uncles
func (pBlock *Block) checkUncleValidity() (isValid bool) {
	return false
}

// Generate Hash of a Block. Using Block header which includes Timestamp, Nonce,
// previous Block Hash
func (pBlock *Block) calculateHash() string {
	bHeader := strconv.Itoa(pBlock.Nonce) + pBlock.Timestamp.String() + pBlock.HashPreviousBlock
	Hash := sha256.New()
	Hash.Write([]byte(bHeader))
	return hex.EncodeToString(Hash.Sum(nil))
}
