package ghost

// *** Structs ***
/* Declaration of structure
Contains blocks and state
*/
type ghost struct {
	blocks []block
	state  map[string]*account
}

// *** Constructors ***
// *** Methods ***

// State transition function. Checks validity of a change in state from a list of transactions
// Syntax APPLY(S,TX) -> S'
func stateTransition(pCurrentState map[string]*account, pTransaction transaction) (pModifiedState map[string]*account, err string) {
	// If referenced UTXO is not in S
	err = ""
	pModifiedState = pCurrentState
	if pCurrentState[pTransaction.origin].balance <= pTransaction.value {
		err = "The referenced UTXO is not in the state\n"
	}
	// If the provided signature does not match the owner of the UTXO
	// TODO: Calculating signature
	//if pTransaction.origin != pTransaction.senderSignature {
	//	err = err + "The provided signature does not match the owner of the UTXO\n"
	//}
	// If the sum of the denominations of all input UTXO is less than the sum of the
	// denominations of all output UTXO, return an error. Not necessary given that a
	// transaction struct only contains one transaction.
	// Creating the account in the state if it doesn't already exist
	if _, ok := pModifiedState[pTransaction.destination]; !ok && err == "" {
		CreateAccount(pTransaction.destination)
	}
	// Return S'. Apply the changes in the transaction
	if err == "" {
		pModifiedState[pTransaction.origin].balance -= pTransaction.value
		pModifiedState[pTransaction.destination].balance += pTransaction.value
	}
	return pModifiedState, err
}

func (pNode *ghostNode) mining() {

}
