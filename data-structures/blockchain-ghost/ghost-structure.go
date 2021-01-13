package ghost

import "errors"

// *** Structs ***
/* Declaration of structure
Contains Blocks and State
*/
type Ghost struct {
	Blocks []Block
	State  map[string]*Account
}

// *** Constructors ***
// *** Methods ***

// State transition function. Checks validity of a change in State from a list of Transactions
// Syntax APPLY(S,TX) -> S'
func stateTransition(pCurrentState map[string]*Account, pTransaction Transaction) (pModifiedState map[string]*Account, err error) {
	// If referenced UTXO is not in S
	pModifiedState = pCurrentState
	if pCurrentState[pTransaction.Origin].Balance <= pTransaction.Value {
		err = errors.New("the referenced UTXO is not in the State")
	}
	// If the provided signature does not match the owner of the UTXO
	// TODO: Calculating signature
	//if pTransaction.Origin != pTransaction.senderSignature {
	//	err = err + "The provided signature does not match the owner of the UTXO\n"
	//}
	// If the sum of the denominations of all input UTXO is less than the sum of the
	// denominations of all output UTXO, return an error. Not necessary given that a
	// Transaction struct only contains one Transaction.
	// Creating the Account in the State if it doesn't already exist
	if _, ok := pModifiedState[pTransaction.Destination]; !ok && err == nil {
		CreateAccount(pTransaction.Destination)
	}
	// Return S'. Apply the changes in the Transaction
	if err == nil {
		pModifiedState[pTransaction.Origin].Balance -= pTransaction.Value
		pModifiedState[pTransaction.Destination].Balance += pTransaction.Value
	}
	return pModifiedState, err
}

// TODO: Mining?
func (pNode *NodeGhost) mining() {

}

// TODO: Handle incoming data streams, to check whether
// When another node connects to our host and wants to propose a new Blockchain to overwrite our own, we need logic to
// determine whether or not we should accept it.

// TODO: Adding new Blocks to the Blockchain and broadcast them
