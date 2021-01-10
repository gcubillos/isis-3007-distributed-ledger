package ghost

// What an account contains
// Nonce counter used to make sure each transaction can only be processed once
// account's current balance
type account struct {
	nonce   int
	balance float32
	address string
}

// *** Constructors ***

// Create an account with balance 0
func CreateAccount(pAddress string) account {
	var rAccount account
	rAccount.balance = 0
	rAccount.address = pAddress
	rAccount.nonce = 0
	return rAccount
}

// *** Methods ***
