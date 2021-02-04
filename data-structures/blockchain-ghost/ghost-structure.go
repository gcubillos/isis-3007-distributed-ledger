package ghost

// *** Structs ***

/* Declaration of structure
Contains Blocks and State, it also saves the children and unused nodes */
type Ghost struct {
	Blocks       []Block
	State        map[string]*Account
	CurrentChain []Block
}

// *** Constructors ***
// *** Methods ***

// Finding the GHOST (Greedy Heaviest-Observed Sub-Tree)
// Way of replacing the chain
// Choosing the branch with the most combined proof of work, measured by the amount of
// nodes present in such branch
// Modified version where you start at the tip and work your way backwards to find the
// heaviest sub tree
func (pGhost *Ghost) FindGHOST(pNewBlockchain Ghost) {
	var forkBlock Block
	var diverges = false
	// Find the place the fork occurs and history diverges
	for i := 0; !diverges; i++ {
		if pGhost.CurrentChain[i].Hash != pNewBlockchain.CurrentChain[i].Hash {
			forkBlock = *pGhost.CurrentChain[i].Parent
			diverges = true
		}
	}

	// Subtree size of the new chain
	newChainSize := findSubTreeSize(pNewBlockchain, forkBlock)
	currentChainSize := findSubTreeSize(*pGhost, forkBlock)

	if newChainSize > currentChainSize {
		pGhost.Blocks = pNewBlockchain.Blocks
		pGhost.CurrentChain = pNewBlockchain.CurrentChain
		pGhost.State = pNewBlockchain.State
	}
}

func findChildren(pBlock Block, pBlockchain Ghost) int {
	numberOfChildren := 0
	for _, v := range pBlockchain.Blocks {
		if v.Parent.Hash == pBlock.Hash {
			numberOfChildren++
		}
	}
	return numberOfChildren
}

func findSubTreeSize(pNewBlockchain Ghost, pForkBlock Block) int {
	pCurrentBlock := pNewBlockchain.CurrentChain[len(pNewBlockchain.CurrentChain)-1]
	newSubtreeSize := 0
	var reachedFork = false
	for j := 0; !reachedFork; j++ {
		if pCurrentBlock.Parent.Hash == pForkBlock.Hash {
			reachedFork = true
		} else {
			newSubtreeSize += 1
			newSubtreeSize += findChildren(pCurrentBlock, pNewBlockchain) - 1
			pCurrentBlock = *pCurrentBlock.Parent
		}
	}
	return newSubtreeSize

}

// TODO: Having two implementations of Ghost, one with the ethereum chain selection and another
// with the one present in the ghost paper
