package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	log2 "log"
	net2 "net"
	"os"
	"strconv"
	"strings"
	"time" 
	"sync"
	// "unsafe"
	"crypto/sha256"
	"encoding/hex"

	"github.com/perlin-network/noise/crypto/ed25519"
	"github.com/perlin-network/noise/examples/chat/messages"
	"github.com/perlin-network/noise/log"
	"github.com/perlin-network/noise/network"
	"github.com/perlin-network/noise/network/discovery"
	"github.com/perlin-network/noise/types/opcode"
)

type ChatPlugin struct{ *network.Plugin }

var mutex = &sync.Mutex{}

func (state *ChatPlugin) Receive(ctx *network.PluginContext) error {
	switch msg := ctx.Message().(type) {
	case *messages.ChatMessage:
		//log.Info().Msgf("<%s> %s", ctx.Client().ID.Address, "Received: "+msg.Message)

		fullMsg := strings.Fields(msg.Message)

		sender := fullMsg[0]
		receiver := fullMsg[1]
		myAmount, err := strconv.Atoi(fullMsg[2])
		if err != nil {
			// handle error
		}

		sendHash := fullMsg[3]
		receiveHash := fullMsg[4]

		//update lattice for sender
		mutex.Lock()
		chain := ctx.Network().Lattice[sender]
		newCube := generateCube(chain[len(chain)-1], "send", myAmount)
		newCube.SendHash = sendHash
		chain = append(chain, newCube)
		ctx.Network().Lattice[sender] = chain
		mutex.Unlock()

		//update lattice for receiver
		mutex.Lock()
		chain = ctx.Network().Lattice[receiver]
		newCube = generateCube(chain[len(chain)-1], "receive", myAmount)
		newCube.ReceiveHash = receiveHash
		chain = append(chain, newCube)
		ctx.Network().Lattice[receiver] = chain
		mutex.Unlock()

		//fmt.Printf("%+v\n", ctx.Network().Lattice)

		// b, err := json.MarshalIndent(ctx.Network().Lattice, "", "  ")
		// if err != nil {
		// 	fmt.Println("error:", err)
		// }
		// fmt.Print(string(b))

		//Latency Test
		fmt.Println("# of transactions: ", len(ctx.Network().Lattice[receiver]))

		timeSentString := fullMsg[5]
		timeSent, err := strconv.ParseInt(timeSentString, 10, 64)
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

	net, err := builder.Build("lattice")
	if err != nil {
		log.Fatal().Err(err)
		return
	}

	go net.Listen()

	if len(peers) > 0 {
		net.Bootstrap(peers...)
	}

	// Tests
	if net.Address == "tcp://192.168.0.22:3000" {

		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')

		//Throughput Test
		// timer := time.NewTimer(time.Second)

		// done := false
		// go func() {
		// 	<-timer.C
		// 	done = true
		// }()

		// for !done {
		// 	mySender := "tcp://10.150.0.2:3000"
		// 	myRecipient := "tcp://10.150.0.4:3000"
		// 	myMsg := "10"
		// 	myAmount, err := strconv.Atoi(myMsg)
		// 	if err != nil {
		// 		// handle error
		// 	}

		// 	//update lattice for sender
		// 	mutex.Lock()
		// 	chain := net.Lattice[mySender]
		// 	newCube := generateCube(chain[len(chain)-1], "send", myAmount)
		// 	newCube.SendHash = calculateHash(newCube)
		// 	sendHash := newCube.SendHash
		// 	chain = append(chain, newCube)
		// 	net.Lattice[mySender] = chain
		// 	mutex.Unlock()

		// 	//update lattice for receiver
		// 	mutex.Lock()
		// 	chain = net.Lattice[myRecipient]
		// 	newCube = generateCube(chain[len(chain)-1], "receive", myAmount)
		// 	newCube.ReceiveHash = calculateHash(newCube)
		// 	receiveHash := newCube.ReceiveHash
		// 	chain = append(chain, newCube)
		// 	net.Lattice[myRecipient] = chain
		// 	mutex.Unlock()

		// 	ctx := network.WithSignMessage(context.Background(), true)
		// 	net.Broadcast(ctx, &messages.ChatMessage{Message: mySender + " " + myRecipient + " " + myMsg + " " + sendHash +" "+ receiveHash + " "+ "1"})
		// }

		// Latency Test
		for i := 0; i < 400; i++ {
			now := time.Now()
			timeNanos := now.UnixNano()

			mySender := "tcp://10.150.0.2:3000"
			myRecipient := "tcp://10.150.0.4:3000"
			myMsg := "10"
			myAmount, err := strconv.Atoi(myMsg)
			if err != nil {
				// handle error
			}

			//update lattice for sender
			mutex.Lock()
			chain := net.Lattice[mySender]
			newCube := generateCube(chain[len(chain)-1], "send", myAmount)
			newCube.SendHash = calculateHash(newCube)
			sendHash := newCube.SendHash
			chain = append(chain, newCube)
			net.Lattice[mySender] = chain
			mutex.Unlock()

			//update lattice for receiver
			mutex.Lock()
			chain = net.Lattice[myRecipient]
			newCube = generateCube(chain[len(chain)-1], "receive", myAmount)
			newCube.ReceiveHash = calculateHash(newCube)
			receiveHash := newCube.ReceiveHash
			chain = append(chain, newCube)
			net.Lattice[myRecipient] = chain
			mutex.Unlock()

			ctx := network.WithSignMessage(context.Background(), true)
			
			timeString := strconv.FormatInt(timeNanos, 10)
			net.Broadcast(ctx, &messages.ChatMessage{Message: mySender + " " + myRecipient + " " + myMsg + " " + sendHash +" "+ receiveHash + " " + timeString})

		}

		// Size Test
		// fmt.Println("Size of Lattice:  ", unsafe.Sizeof(net.Lattice))
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')

		// skip blank lines
		if len(strings.TrimSpace(input)) == 0 {
			continue
		}

		fullMsg := strings.Fields(input)

		myRecipient := fullMsg[0]
		myMsg := fullMsg[1]
		myAmount, err := strconv.Atoi(myMsg)
		if err != nil {
			// handle error
		}

		//update lattice for sender
		chain := net.Lattice[net.Address]
		newCube := generateCube(chain[len(chain)-1], "send", myAmount)
		newCube.SendHash = calculateHash(newCube)
		sendHash := newCube.SendHash
		chain = append(chain, newCube)
		net.Lattice[net.Address] = chain

		//update lattice for receiver
		chain = net.Lattice[myRecipient]
		newCube = generateCube(chain[len(chain)-1], "receive", myAmount)
		newCube.ReceiveHash = calculateHash(newCube)
		receiveHash := newCube.ReceiveHash
		chain = append(chain, newCube)
		net.Lattice[myRecipient] = chain

		ctx := network.WithSignMessage(context.Background(), true)
		net.Broadcast(ctx, &messages.ChatMessage{Message: net.Address + " " + myRecipient + " " + myMsg + " " + sendHash + " " +receiveHash + " "+ "1"})

		//fmt.Printf("%+v\n", net.Lattice)

		b, err := json.MarshalIndent(net.Lattice, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Print(string(b))
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

// create a new cube
func generateCube(oldCube network.Cube, typeOfTransaction string, amount int) network.Cube {

	var newCube network.Cube

	newCube.Index = oldCube.Index + 1

	if typeOfTransaction == "send" {
		newCube.Balance = oldCube.Balance - amount
	} else {
		newCube.Balance = oldCube.Balance + amount
	}

	newCube.Type = typeOfTransaction
	newCube.Amount = amount

	return newCube
}

// SHA256 hashing
func calculateHash(cube network.Cube) string {
	// record := strconv.Itoa(block.Index) + block.Timestamp +
	// 	strconv.Itoa(block.BPM) + block.PrevHash + block.Nonce
	record := strconv.Itoa(cube.Amount + cube.Balance + cube.Index)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}
