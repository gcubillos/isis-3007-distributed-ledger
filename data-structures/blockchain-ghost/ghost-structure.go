package ghost

// *** Structs ***

/* Declaration of structure
Contains Blocks and State, it also saves the children and unused nodes */
type Ghost struct {
	Blocks      []Block
	State       map[string]*Account
	CurrentLeaf Block
}

// *** Constructors ***
// *** Methods ***

// Way of replacing the chain
// Choosing the branch with the most combined proof of work, measured by the amount of
// nodes present in such branch
// TODO: Finish replace chain algorithm

// Finding the GHOST (Greedy Heaviest-Observed Sub-Tree)
// TODO: Finish finding GHOST
func (pGhost *Ghost) FindGHOST(pNewBlockchain Ghost) {
	// Find the forking and divergence in history

	// Start in genesis block
	currentBlock := pGhost.Blocks[0]
	// Find children of block
	newBlockChildren := make([]Block, 0)
	for _, v := range pNewBlockchain.Blocks {
		if v.Parent.Hash == currentBlock.Hash {
			newBlockChildren = append(newBlockChildren, v)
		}
	}
	// Check result
	if len(newBlockChildren) == 0 {
		// Update the leaf
		pGhost.CurrentLeaf
	} else {
		biggestSubtree := 0
		for _, v := range newBlockChildren {
			if biggestSubtree < v.Children {

			}

		}
	}

}

// TODO: Having two implementations of Ghost, one with the ethereum chain selection and another
// with the one present in the ghost paper
