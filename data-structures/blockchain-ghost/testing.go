package main

import "github.com/perlin-network/noise"


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

// When Bob gets a message from Alice, print it out and respond to Alice with 'Hi Alice!'

bob.Handle(func(ctx noise.HandlerContext) error {
if !ctx.IsRequest() {
return nil
}

fmt.Printf("Got a message from Alice: '%s'\n", string(ctx.Data()))

return ctx.Send([]byte("Hi Alice!"))
})

// Have Alice and Bob start listening for new peers.

if err := alice.Listen(); err != nil {
panic(err)
}

if err := bob.Listen(); err != nil {
panic(err)
}

// Have Alice send Bob a request with the message 'Hi Bob!'

res, err := alice.Request(context.TODO(), bob.Addr(), []byte("Hi Bob!"))
if err != nil {
panic(err)
}

// Print out the response Bob got from Alice.

fmt.Printf("Got a message from Bob: '%s'\n", string(res))