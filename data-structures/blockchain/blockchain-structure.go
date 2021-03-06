package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gcubillos/isis-3007-distributed-ledger/data-structures/shared-components"
	"strconv"
	"strings"
	"time"
)

// *** Structs ***

// TODO: Way to define the difficulty in the caller program
// The number of leading zeroes wanted from the hash when doing the proof of work
var Difficulty = 1

// What a block in the blockchain contains
type Block struct {
	Timestamp    time.Time
	Hash         string
	PrevHash     string
	Nonce        int
	Transactions []components.Transaction
	Difficulty   int
}

// What the blockchain data structure contains
type Blockchain struct {
	Blocks []Block
	State  map[string]float64
}

// *** Methods ***

// Generate Hash of a block
func CalculateHash(block Block) string {
	record := strconv.Itoa(block.Nonce) + block.Timestamp.String() + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}

// Function that checks whether a block is valid
// TODO: Receives a boolean parameter that indicates whether it is a new block or not.
// This is done as a way to reverse the state and checking it
func (pBlockchain *Blockchain) IsBlockValid(newBlock, oldBlock Block) (bool, error) {
	// Validating previous blocks until the genesis block is reached
	if oldBlock.PrevHash != "" {
		// Finding parent of block
		var parentOfBlock Block
		for _, v := range pBlockchain.Blocks {
			if v.PrevHash == oldBlock.Hash {
				parentOfBlock = v
				break
			}
		}
		// Previous block exists and is valid
		if condition, _ := pBlockchain.IsBlockValid(oldBlock, parentOfBlock); !condition {
			return false, errors.New("previous Block isn't valid")
		}
		switch true {
		// Previous block exists in the blockchain. It is assumed it is valid
		case oldBlock.Hash != pBlockchain.Blocks[len(pBlockchain.Blocks)-1].Hash:
			return false, errors.New("oldBlock's hash doesn't seem to match the latest block in the blockchain")
		// Timestamp
		case !oldBlock.Timestamp.Before(newBlock.Timestamp):
			return false, errors.New("timestamp is not valid")
		// Previous block hash comparison
		case oldBlock.Hash != newBlock.PrevHash:
			return false, errors.New("hash of previous block doesn't match")
		// Does the corresponding hash match
		case CalculateHash(newBlock) != newBlock.Hash:
			return false, errors.New("calculated hash doesn't match")
		// Checking proof of work
		case !IsHashValid(newBlock.Hash, newBlock.Difficulty):
			return false, errors.New("the proof of work is not valid")
		// Verifying state transition
		case !pBlockchain.verifyStateTransition(newBlock.Transactions, pBlockchain.State):
			return false, errors.New("the transactions are inconsistent with the state")
		default:
			return true, nil
		}

	} else {
		return true, nil
	}
}

// Checks whether the hash is valid by checking if it starts with the given number of zeroes specified in the difficulty
// TODO: Check whether it influences if the block has more than the difficulty number of leading zeroes. Does it matter?
func IsHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

// that is that the transitions in the state are valid
func (pBlockchain *Blockchain) ReplaceChain(newBlockchain Blockchain) {
	var forkBlock Block
	diverges := false
	// Find the place the fork occurs and history diverges
	for i := 0; !diverges && i < len(pBlockchain.Blocks); i++ {
		if pBlockchain.Blocks[i].Hash != newBlockchain.Blocks[i].Hash {
			forkBlock = pBlockchain.Blocks[i-1]
			diverges = true
		}
	}
	if forkBlock.Hash != pBlockchain.Blocks[0].Hash {
		lenNewBlockchain := len(newBlockchain.Blocks)
		// Verify length and validity of blocks in the received chain
		if ok, _ := newBlockchain.IsBlockValid(newBlockchain.Blocks[lenNewBlockchain-1], newBlockchain.Blocks[lenNewBlockchain-2]); lenNewBlockchain > len(pBlockchain.Blocks) && ok {
			pBlockchain.Blocks = newBlockchain.Blocks
			pBlockchain.State = newBlockchain.State
		}
	}
}

// Receives a state and then performs the transactions and returns the modified state when it is valid
func (pBlockchain *Blockchain) verifyStateTransition(pTransactions []components.Transaction, initialState map[string]float64) bool {
	modifiedState := initialState
	for _, v := range pTransactions {
		switch true {
		// TODO: Verifying signature, doing it in the same main function?
		// Signature of sender does not match the owner of the UTXO
		// Transaction is well formed
		case pBlockchain.State[v.Origin] < 0:
			return false
		// UTXO is not in the state
		case pBlockchain.State[v.Origin] < v.Value:
			return false
		}
		// Update state
		modifiedState[v.Origin] -= v.Value
		// Checking that the recipient of the UTXO exists. If not, create it
		if _, ok := modifiedState[v.Destination]; ok {
			modifiedState[v.Destination] += v.Value
		} else {
			modifiedState[v.Destination] = v.Value
		}
	}
	// Update the final state
	initialState = modifiedState
	return true
}

// TODO: Adding a limit for number of transactions in a block?
// TODO: Add a function to request to add a transaction to a block
// TODO: Standardize names through out the implementations
