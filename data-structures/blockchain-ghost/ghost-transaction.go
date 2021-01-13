package ghost

// What a transaction ensues
// A transaction is a request to move $X from A to B
type Transaction struct {
	Origin          string
	SenderSignature string
	Destination     string
	Value           float64
}
