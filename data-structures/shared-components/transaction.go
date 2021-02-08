package components

// What a transaction ensues
// A transaction is a request to move $X from A to B
type Transaction struct {
	Origin          string
	SenderSignature string
	Destination     string
	Value           float64
}

func CreateTransaction(pOrigin, pSignature, pDestination string, pValue float64) Transaction {
	return Transaction{
		Origin:          pOrigin,
		SenderSignature: pSignature,
		Destination:     pDestination,
		Value:           pValue,
	}
}
