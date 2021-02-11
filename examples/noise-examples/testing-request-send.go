package main

import (
	"context"
	"fmt"
	"github.com/perlin-network/noise"
)

// This example demonstrates how to send/handle RPC requests across peers, how to listen for incoming
// peers, how to check if a message received is a request or not, how to reply to a RPC request, and
// how to cleanup node instances after you are done using them.
func main() {
	// Let there be nodes Alice and Bob.

	alice, err := noise.NewNode()
	check(err)

	bob, err := noise.NewNode()
	check(err)

	charlie, err := noise.NewNode()
	check(err)

	//// Gracefully release resources for Alice and Bob at the end of the example.
	//
	//defer alice.Close()
	//defer bob.Close()

	// When Bob gets a message from Alice, print it out and respond to Alice with 'Hi Alice!'

	bob.Handle(func(ctx noise.HandlerContext) error {

		fmt.Printf("Got a message from Alice: '%s'\n", string(ctx.Data()))

		return ctx.Send([]byte("Hi Alice!"))
	})

	charlie.Handle(func(ctx noise.HandlerContext) error {

		fmt.Printf("Got a message from Alice: '%s'\n", string(ctx.Data()))

		return ctx.Send([]byte("Hi Alice!"))
	})

	alice.Handle(func(ctx noise.HandlerContext) error {

		for _, v := range alice.Outbound() {
			_, err := alice.Request(context.TODO(), v.ID().Address, []byte("Second message"))
			check(err)
		}

		return ctx.Send([]byte(""))
	})

	// Have Alice and Bob start listening for new peers.

	check(alice.Listen())
	check(bob.Listen())
	check(charlie.Listen())

	// Have Alice send Bob a request with the message 'Hi Bob!'

	res, err := alice.Request(context.TODO(), bob.Addr(), []byte("Hi Bob!"))
	check(err)

	res2, err := alice.Request(context.TODO(), charlie.Addr(), []byte("Hi Charlie!"))
	check(err)

	res3, err := charlie.Request(context.TODO(), alice.Addr(), []byte("Testing Alice!"))

	// Print out the response Bob got from Alice.

	fmt.Printf("Got a message from Bob: '%s'\n", string(res))

	fmt.Printf("Another message: '%s'\n", string(res2))

	fmt.Printf("res3 '%s'\n", string(res3))

	// Output:
	// Got a message from Alice: 'Hi Bob!'
	// Got a message from Bob: 'Hi Alice!'
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
