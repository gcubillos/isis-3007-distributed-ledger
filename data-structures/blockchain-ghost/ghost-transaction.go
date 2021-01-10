package ghost

// What a transaction ensues
// A transaction is a request to move $X from A to B
type transaction struct {
	origin          string
	senderSignature string
	destination     string
	value           float32
}
