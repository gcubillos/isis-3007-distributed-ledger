package main

import (
	"fmt"
)

// Declaration of structure

type ghost struct {
	blocks []block
}


// Declaration of a block in the network
// The network is intended to produce roughly one block every ten minutes, with each block 
// containing a timestamp, a nonce, a reference to (ie. hash of) the previous block and a list of all of the transactions that have taken place since the previous block.

type block struct {
	timestamp int
	nonce int
	hashPreviousBlock string
	hashNextBlock string

}

// Check that the block is valid
// By checking if the previous block referenced by the block exists and is valid
// Checking that the timestamp of the block is greater than that of the previous block
// Check that the proof of work on the block is valid.
// Let S[0] be the state at the end of the previous block.
// Suppose TX is the block's transaction list with n transactions. For all i in 0...n-1, set S[i+1] = APPLY(S[i],TX[i]) If any application returns an error, exit and return false.
// Return true, and register S[n] as the state at the end of this block.

func checkBlockValid(pBlock block) (z bool) {
	z = false
	if pBlock.nonce < 1 {
		z = true
	}
	return
}

func main(){
	fmt.Printf("Hello, trying to build ghost")
	var theGhost = make(ghost,0)
	var theBlock = make(block,0)
	theBlock.nonce = 2
	theGhost.blocks = append(theGhost.blocks, theBlock)

}
