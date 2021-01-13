package ghost

// What an account contains
// Nonce counter used to make sure each transaction can only be processed once
// account's current Balance
type Account struct {
	Nonce   int
	Balance float64
	Address string
}

// *** Constructors ***

// Create an account with Balance 0
func CreateAccount(pAddress string) Account {
	var rAccount Account
	rAccount.Balance = 0
	rAccount.Address = pAddress
	rAccount.Nonce = 0
	return rAccount
}

// *** Methods ***
