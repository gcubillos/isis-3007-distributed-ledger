package ghost

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gcubillos/isis-3007-distributed-ledger/data-structures/shared-components"
	"strconv"
	"strings"
	"time"
)

// What a Block in the network contains
// The network is intended to produce roughly one Block every ten minutes, with each Block
// containing a Timestamp, a Nonce, a reference to (ie. Hash of) the previous Block and a
// list of all of the Transactions that have taken place since the previous Block.

type Block struct {
	Timestamp         time.Time
	Nonce             int
	Hash              string
	HashPreviousBlock string
	Parent            *Block
	Uncles            []Block
	Transactions      []components.Transaction
	RecentState       map[string]*Account
	BlockNumber       int
	Difficulty        int
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

func (pGhost *Ghost) IsBlockValid(pBlock Block) (bool, error) {
	// Previous Block exists and valid
	if condition, _ := pGhost.IsBlockValid(*pBlock.Parent); !condition {
		return false, errors.New("previous Block isn't valid")
	}
	switch true {
	// Timestamp
	case pBlock.Timestamp.Before(pBlock.Parent.Timestamp):
		return false, errors.New("timestamp of previous Block isn't valid")
	// Previous block hash comparison
	case pBlock.HashPreviousBlock != calculateHash(pBlock):
		return false, errors.New("hash of previous block doesn't match")
	// Checking that the current hash is valid
	case pBlock.Hash != calculateHash(pBlock):
		return false, errors.New("current block hash is not valid")
	// Validating proof of work
	case !IsHashValid(pBlock.Hash, pBlock.Difficulty):
		return false, errors.New("Proof of work is not valid")
	// State transition check
	// TODO: Add the state if necessary, maybe it can be accessed through the block itself
	case !pGhost.verifyStateTransition(pBlock.Transactions, nil):
		return false, errors.New("the transactions are inconsistent with the state")
	default:
		return true, nil
	}
}

//Checks validity of uncles
//TODO: Finish validity of uncles
func (pBlock *Block) checkUncleValidity() (isValid bool) {
	return false
}

// Generate Hash of a Block. Using Block header which includes Timestamp, Nonce,
// previous Block Hash
func calculateHash(pBlock Block) string {
	bHeader := strconv.Itoa(pBlock.Nonce) + pBlock.Timestamp.String() + pBlock.HashPreviousBlock
	Hash := sha256.New()
	Hash.Write([]byte(bHeader))
	return hex.EncodeToString(Hash.Sum(nil))
}

// Checks whether the hash is valid by checking if it starts with the given number of zeroes specified in the difficulty
// TODO: Check whether it influences if the block has more than the difficulty number of leading zeroes. Does it matter?
func IsHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

// Receives a state and then performs the transactions and returns the modified state when it is valid
func (pGhost *Ghost) verifyStateTransition(pTransactions []components.Transaction, initialState map[string]float64) bool {
case:

	var initialState = pBlock.Parent.RecentState
	for i := 0; i < len(pBlock.Transactions); i++ {
	if currentState, err := stateTransition(initialState, pBlock.Transactions[i]); err != nil {
	initialState = currentState
	} else {
	break
	}
	}
	// Checking state
	if len(pBlock.RecentState) == len(initialState) {
	for i := range pBlock.RecentState {
	if !(pBlock.RecentState[i].Nonce == initialState[i].Nonce) &&
	(pBlock.RecentState[i].Balance == initialState[i].Balance) &&
	(pBlock.RecentState[i].Address == initialState[i].Address) {
	return false, errors.New("there is an error with the state")
	}
	}
	} else {
	return false, errors.New("state doesn't match")
	}
	modifiedState := initialState
	for _, v := range pTransactions {
		// TODO: Verifying signature, doing it in the same main function?
		// TODO: Change it so that it verifies the signature
		// Signature of sender does not match the owner of the UTXO
		// UTXO is not in the state
		if pGhost.State[v.Origin] < v.Value {
			return false
		} else {
			// Update state
			modifiedState[v.Origin] -= v.Value
			// Checking that the recipient of the UTXO exists. If not, create it
			if _, ok := modifiedState[v.Destination]; ok {
				modifiedState[v.Destination] += v.Value
			} else {
				modifiedState[v.Destination] = v.Value
			}
		}
	}
	// Update the final state
	initialState = modifiedState
	return true
	// TODO: Adding a limit for number of transactions in a block?
	// TODO: Add a function to request to add a transaction to a block
}