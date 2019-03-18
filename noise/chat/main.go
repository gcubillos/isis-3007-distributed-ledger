package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	log2 "log" //"math/big"
	net2 "net"
	"os"
	"strconv"
	"strings"
	"time" //"unsafe"

	"github.com/perlin-network/noise/crypto/ed25519"
	"github.com/perlin-network/noise/examples/chat/messages"
	"github.com/perlin-network/noise/log"
	"github.com/perlin-network/noise/network"
	"github.com/perlin-network/noise/network/discovery"
	"github.com/perlin-network/noise/types/opcode"
)

type ChatPlugin struct{ *network.Plugin }

func (state *ChatPlugin) Receive(ctx *network.PluginContext) error {
	switch msg := ctx.Message().(type) {
	case *messages.ChatMessage:
		log.Info().Msgf("<%s> %s", ctx.Client().ID.Address, msg.Message)

		//Latency Test
		timeSent, err := strconv.ParseInt(msg.Message, 10, 64)
		if err != nil {
			fmt.Println(err)
		}

		now := time.Now()
		timeNanos := now.UnixNano()

		nanos := timeNanos - timeSent
		fmt.Printf("Latency: %dns", nanos)
		fmt.Println()
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

	net, err := builder.Build("chat")
	if err != nil {
		log.Fatal().Err(err)
		return
	}

	go net.Listen()

	if len(peers) > 0 {
		net.Bootstrap(peers...)
	}

	// For Tests
	if net.Address == "tcp://192.168.0.16:3000" {

		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')

		// // Throughput Test
		// timer := time.NewTimer(time.Second)

		// done := false
		// go func() {
		// 	<-timer.C
		// 	done = true
		// }()

		// i := 0
		// for !done {
		// 	myMessage := "message " + strconv.Itoa(i)
		// 	i++
		// 	log.Info().Msgf("<%s> %s", net.Address, myMessage)
		// 	ctx := network.WithSignMessage(context.Background(), true)
		// 	net.Broadcast(ctx, &messages.ChatMessage{Message: myMessage})
		// }

		// Latency Test
		for i := 0; i < 400; i++ {
			now := time.Now()
			timeNanos := now.UnixNano()

			myMessage := strconv.FormatInt(timeNanos, 10)
			log.Info().Msgf("<%s> %s", net.Address, myMessage)
			ctx := network.WithSignMessage(context.Background(), true)
			net.Broadcast(ctx, &messages.ChatMessage{Message: myMessage})

			// Size Test
			// 	myMessage := "single message"
			// 	log.Info().Msgf("<%s> %s", net.Address, myMessage)
			// 	ctx := network.WithSignMessage(context.Background(), true)
			// 	net.Broadcast(ctx, &messages.ChatMessage{Message: myMessage})
			// 	fmt.Println("Size of ChatMessage: ", unsafe.Sizeof(myMessage))
		}
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')

		// skip blank lines
		if len(strings.TrimSpace(input)) == 0 {
			continue
		}

		log.Info().Msgf("<%s> %s", net.Address, input)
		ctx := network.WithSignMessage(context.Background(), true)
		net.Broadcast(ctx, &messages.ChatMessage{Message: input})
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
