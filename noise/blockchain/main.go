package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	log2 "log"
	net2 "net"
	"os"
	"strconv"
	"strings"

	"github.com/perlin-network/noise/crypto/ed25519"
	"github.com/perlin-network/noise/examples/chat/messages"
	"github.com/perlin-network/noise/log"
	"github.com/perlin-network/noise/network"
	"github.com/perlin-network/noise/network/discovery"
	"github.com/perlin-network/noise/types/opcode"
)

type ChatPlugin struct{ *network.Plugin }

type block struct {
	Index  int
	Type   string
	Amount string
}

func (state *ChatPlugin) Receive(ctx *network.PluginContext) error {
	switch msg := ctx.Message().(type) {
	case *messages.ChatMessage:
		log.Info().Msgf("<%s> %s", ctx.Client().ID.Address, "Received: "+msg.Message)

		myAmount, err := strconv.Atoi(msg.Message)
		if err != nil {
			// handle error
		}

		//update blockchain
		newBlock := generateBlock(ctx.Network().Blockchain[len(ctx.Network().Blockchain)-1], "receive", myAmount)
		ctx.Network().Blockchain = append(ctx.Network().Blockchain, newBlock)

		fmt.Printf("%+v\n", ctx.Network().Blockchain)
	}

	return nil
}

func main() {
	// process other flags
	portFlag := flag.Int("port", 3000, "port to listen to")
	hostFlag := flag.String("host", getOutboundIP(), "host to listen to")
	protocolFlag := flag.String("protocol", "tcp", "protocol to use (kcp/tcp)")
	peersFlag := flag.String("peers", "", "peers to connect to")
	flag.Parse()

	port := uint16(*portFlag)
	host := *hostFlag
	protocol := *protocolFlag
	peers := strings.Split(*peersFlag, ",")

	keys := ed25519.RandomKeyPair()

	// log.Info().Msgf("Private Key: %s", keys.PrivateKeyHex())
	// log.Info().Msgf("Public Key: %s", keys.PublicKeyHex())

	opcode.RegisterMessageType(opcode.Opcode(1000), &messages.ChatMessage{})
	builder := network.NewBuilder()
	builder.SetKeys(keys)
	builder.SetAddress(network.FormatAddress(protocol, host, port))

	// Register peer discovery plugin.
	builder.AddPlugin(new(discovery.Plugin))

	// Add custom chat plugin.
	builder.AddPlugin(new(ChatPlugin))

	net, err := builder.Build("blockchain")
	if err != nil {
		log.Fatal().Err(err)
		return
	}

	go net.Listen()

	if len(peers) > 0 {
		net.Bootstrap(peers...)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')

		// skip blank lines
		if len(strings.TrimSpace(input)) == 0 {
			continue
		}

		ss := strings.Fields(input)

		myRecipient := ss[0]
		myMsg := ss[1]
		myAmount, err := strconv.Atoi(myMsg)
		if err != nil {
			// handle error
		}

		ctx := network.WithSignMessage(context.Background(), true)

		if client, err := net.Client(myRecipient); err == nil {
			client.Tell(ctx, &messages.ChatMessage{Message: myMsg})
			log.Info().Msgf("<%s> %s", net.Address, "Sent: "+myMsg)

			//update blockchain
			newBlock := generateBlock(net.Blockchain[len(net.Blockchain)-1], "send", myAmount)
			net.Blockchain = append(net.Blockchain, newBlock)

			fmt.Printf("%+v\n", net.Blockchain)
		}
	}

}

func getOutboundIP() string {
	conn, err := net2.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log2.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net2.UDPAddr)

	return localAddr.IP.String()
}

// create a new block using previous block's hash
func generateBlock(oldBlock network.Block, typeOfTransaction string, amount int) network.Block {

	var newBlock network.Block

	newBlock.Index = oldBlock.Index + 1

	if typeOfTransaction == "send" {
		newBlock.Balance = oldBlock.Balance - amount
	} else {
		newBlock.Balance = oldBlock.Balance + amount
	}

	newBlock.Type = typeOfTransaction
	newBlock.Amount = amount

	return newBlock
}
