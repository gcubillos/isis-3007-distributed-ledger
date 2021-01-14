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
var difficulty = 1

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

// Variable containing the current version of the blockchain
var CurrentBlockchain Blockchain

// *** Methods ***

// Generate Hash of a block
func CalculateHash(block Block) string {
	record := strconv.Itoa(block.Nonce) + block.Timestamp.String() + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Creates a block
func GenerateBlock(oldBlock Block, pTransactions []ghost.Transaction) Block {

	var newBlock Block

	newBlock.Timestamp = time.Now()
	newBlock.Transactions = pTransactions
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty
	// Calculating the hash
	for i := 0; ; i++ {
		newBlock.Nonce = i
		if !IsHashValid(CalculateHash(newBlock), newBlock.Difficulty) {
			continue
		} else {
			newBlock.Hash = CalculateHash(newBlock)
			break
		}
	}

	return newBlock
}

// Function that checks whether a block is valid
func IsBlockValid(newBlock, oldBlock Block) (bool, error) {
	// TODO: Previous block exists and is valid
	// Timestamp
	if !oldBlock.Timestamp.Before(newBlock.Timestamp) {
		return false, errors.New("timestamp is not valid")
	}
	// Previous block hash comparison
	if oldBlock.Hash != newBlock.PrevHash {
		return false, errors.New("hash of previous block doesn't match")
	}
	// Does the corresponding hash match
	if CalculateHash(newBlock) != newBlock.Hash {
		return false, errors.New("calculated hash doesn't match")
	}
	// Checking proof of work
	if !IsHashValid(newBlock.Hash, newBlock.Difficulty) {
		return false, errors.New("the proof of work is not valid")
	}
	if !verifyStateTransition(newBlock.Transactions, CurrentBlockchain.State) {
		return false, errors.New("the transactions are inconsistent with the state")
	}

	return true, nil
}

// Checks whether the hash is valid by checking if it starts with the given number of zeroes specified in the difficulty
// TODO: Check whether it influences if the block has more than the difficulty number of leading zeroes. Does it matter?
func IsHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}

// TODO: Change consensus algorithm?
func ReplaceChain(newBlockchain Blockchain) {
	if len(newBlockchain.Blocks) > len(CurrentBlockchain.Blocks) {
		CurrentBlockchain.Blocks = newBlockchain.Blocks
	}
}

// TODO: Manage concurrency

// Receives a state and then performs the transactions and modifies the given state
func stateTransition(pTransactions []ghost.Transaction, initialState map[string]float64) {
	for _, v := range pTransactions {
		// Update state
		initialState[v.Origin] -= v.Value
		// Checking that the recipient of the UTXO exists. If not, create it
		if _, ok := initialState[v.Destination]; ok {
			initialState[v.Destination] += v.Value
		} else {
			initialState[v.Destination] = v.Value
		}
	}

}

// Receives a state and then performs the transactions and returns the modified state when it is valid
func verifyStateTransition(pTransactions []ghost.Transaction, initialState map[string]float64) bool {
	modifiedState := initialState
	for _, v := range pTransactions {
		// TODO: Verifying signature, doing it in the same main function?
		// TODO: Change it so that it verifies the signature
		// Signature of sender does not match the owner of the UTXO
		// UTXO is not in the state
		if CurrentBlockchain.State[v.Origin] < v.Value {
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
