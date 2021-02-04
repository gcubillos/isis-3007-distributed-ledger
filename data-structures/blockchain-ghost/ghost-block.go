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
	// Checking that it is valid until you reach the genesis block
	if pBlock.HashPreviousBlock != "" {
		// Previous Block exists and valid
		if condition, _ := pGhost.IsBlockValid(*pBlock.Parent); !condition {
			return false, errors.New("previous Block isn't valid")
		}
		switch true {
		// Timestamp
		case pBlock.Timestamp.Before(pBlock.Parent.Timestamp):
			return false, errors.New("timestamp of previous Block isn't valid")
		// Previous block hash comparison
		case pBlock.HashPreviousBlock != CalculateHash(pBlock):
			return false, errors.New("hash of previous block doesn't match")
		// Checking that the current hash is valid
		case pBlock.Hash != CalculateHash(pBlock):
			return false, errors.New("current block hash is not valid")
		// Validating proof of work
		case !IsHashValid(pBlock.Hash, pBlock.Difficulty):
			return false, errors.New("proof of work is not valid")
		// State transition check
		case !verifyStateTransition(pBlock):
			return false, errors.New("the transactions are inconsistent with the state")
		default:
			return true, nil
		}
	} else {
		return true, nil
	}
}

// Generate Hash of a Block. Using Block header which includes Timestamp, Nonce,
// previous Block Hash
func CalculateHash(pBlock Block) string {
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
func verifyStateTransition(pBlock Block) bool {
	// TODO: Checking validity of accounts
	// Initialize state
	var modifiedState = pBlock.Parent.RecentState
	// Go through the lists of transactions
	for _, v := range pBlock.Transactions {
		switch true {
		// Checking transaction is valid and well formed
		case v.Value < 0:
			return false
		// Signature of sender does not match owner
		// TODO: Calculating signature
		// Referenced UTXO is not in the state
		case modifiedState[v.Origin].Balance < v.Value:
			return false
		}
		// Update state
		modifiedState[v.Origin].Balance -= v.Value
		// Check that the recipient of the UTXO exists, if not, create it
		if _, ok := modifiedState[v.Destination]; ok {
			modifiedState[v.Destination].Balance += v.Value
		} else {
			theAccount := CreateAccount(v.Destination)
			modifiedState[v.Destination] = &theAccount
			modifiedState[v.Destination].Balance = v.Value
		}
	}
	// Update the state
	pBlock.RecentState = modifiedState
	return true
}
