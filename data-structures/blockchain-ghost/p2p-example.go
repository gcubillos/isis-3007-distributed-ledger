package main

import (
	"context"
	"fmt"
	"github.com/perlin-network/noise"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type aNetworkNode struct {
	theNode *noise.Node
}

func (p *aNetworkNode) handle() {
	p.theNode.Handle(func(ctx noise.HandlerContext) error {
		if !ctx.IsRequest() {
			return nil
		}

		fmt.Printf("The mega node is activated, Hello World: '%s'\n", string(ctx.Data()))

		return ctx.Send([]byte("Hi Alice!"))
	})

}

func (p *aNetworkNode) listen() error {
	return p.theNode.Listen()
}

// This example demonstrates how to send/handle RPC requests across peers, how to listen for incoming
// peers, how to check if a message received is a request or not, how to reply to a RPC request, and
// how to cleanup node instances after you are done using them.
func main() {
	// Let there be nodes Alice and Bob.

	alice, err := noise.NewNode()
	check(err)

	bob, err := noise.NewNode()
	check(err)

	claire, err := noise.NewNode()
	check(err)

	// Creating a super-node
	superclaire := aNetworkNode{claire}

	// Gracefully release resources for Alice and Bob at the end of the example.
	defer alice.Close()
	defer bob.Close()

	// When Bob gets a message from Alice, print it out and respond to Alice with 'Hi Alice!'
	// Have Alice and Bob start listening for new peers.

	check(alice.Listen())
	check(bob.Listen())
	check(superclaire.listen())

	// Have Alice send Bob a request with the message 'Hi Bob!'

	res, err := alice.Request(context.TODO(), claire.Addr(), []byte("Hi Claire!"))
	check(err)

	// Print out the response Bob got from Alice.

	fmt.Printf("Got a message from Bob: '%s'\n", string(res))

	// Print out address of nodes
	fmt.Printf("Address 1:  %s \n Address 2: %s", alice.Addr(), bob.Addr())

	// Output:
	// Got a message from Alice: 'Hi Bob!'
	// Got a message from Bob: 'Hi Alice!'

}
