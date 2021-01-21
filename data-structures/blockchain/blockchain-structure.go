package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	ghost "github.com/gcubillos/isis-3007-distributed-ledger/data-structures/blockchain-ghost"
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
	Transactions []ghost.Transaction
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
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Function that checks whether a block is valid
func (pBlockchain *Blockchain) IsBlockValid(newBlock, oldBlock Block) (bool, error) {
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
	case !pBlockchain.verifyStateTransition(newBlock.Transactions, pBlockchain.State):
		return false, errors.New("the transactions are inconsistent with the state")
	default:
		return true, nil
	}
}

// Checks whether the hash is valid by checking if it starts with the given number of zeroes specified in the difficulty
// TODO: Check whether it influences if the block has more than the difficulty number of leading zeroes. Does it matter?
func IsHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

// TODO: Change consensus algorithm?
func (pBlockchain *Blockchain) ReplaceChain(newBlockchain Blockchain) {
	if len(newBlockchain.Blocks) > len(pBlockchain.Blocks) {
		pBlockchain.Blocks = newBlockchain.Blocks
	}
}

// TODO: Manage concurrency
// Receives a state and then performs the transactions and returns the modified state when it is valid
func (pBlockchain *Blockchain) verifyStateTransition(pTransactions []ghost.Transaction, initialState map[string]float64) bool {
	modifiedState := initialState
	for _, v := range pTransactions {
		// TODO: Verifying signature, doing it in the same main function?
		// TODO: Change it so that it verifies the signature
		// Signature of sender does not match the owner of the UTXO
		// UTXO is not in the state
		if pBlockchain.State[v.Origin] < v.Value {
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

}
