package main

import (
	"context"
	"fmt"
	"github.com/perlin-network/noise"
)

func main() {
	// Let there be nodes Alice and Bob.

	alice, err := noise.NewNode()
	if err != nil {
		panic(err)
	}

	bob, err := noise.NewNode()
	if err != nil {
		panic(err)
	}

	// Gracefully release resources for Alice and Bob at the end of the example.

	defer alice.Close()
	defer bob.Close()

	// When Alice gets a message from Bob, print it out.

	alice.Handle(func(ctx noise.HandlerContext) error {
		fmt.Printf("Got a message from Bob: '%s'\n", string(ctx.Data()))
		return nil
	})

	// When Bob gets a message from Alice, print it out.

	bob.Handle(func(ctx noise.HandlerContext) error {
		fmt.Printf("Got a message from Alice: '%s'\n", string(ctx.Data()))
		return nil
	})

	// Have Alice and Bob start listening for new peers.

	if err := alice.Listen(); err != nil {
		panic(err)
	}

	if err := bob.Listen(); err != nil {
		panic(err)
	}

	// Have Alice send Bob 'Hi Bob!'

	if err := alice.Send(context.TODO(), bob.Addr(), []byte("Hi Bob!")); err != nil {
		panic(err)
	}

	// Wait until Bob receives the message from Alice.

	// Have Bob send Alice 'Hi Alice!'

	if err := bob.Send(context.TODO(), alice.Addr(), []byte("Hi Alice!")); err != nil {
		panic(err)
	}

	// Wait until Alice receives the message from Bob.

}
