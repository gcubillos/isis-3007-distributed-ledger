package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	mrand "math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	golog "github.com/ipfs/go-log"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	network "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
	gologging "github.com/whyrusleeping/go-logging"
)

type Transaction struct {
	Index      int
	Operation  string
	TimeInt    int64
	TimeString string
	Weight     int
	CumWeight  int
}

type Link struct {
	Target int
	Source int
}

type tangle struct {
	Transactions []Transaction
	Links        []Link
	State        map[string]int
}

// Tangle is a DAG of Transactions
var Tangle struct {
	Transactions []Transaction
	Links        []Link
	lambda       float32
	alpha        float32
	h            int64
	tipSelection string
	State        map[string]int
}

var mutex = &sync.Mutex{}

//Metrics
var initialTime int64
var throughput int64
var latency int64
var size int

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
func makeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	//myIP4 := getOutboundIP()

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		//libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", myIP4, listenPort)),
		libp2p.Identity(priv),
	}

	basicHost, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	log.Printf("I am %s\n", fullAddr)
	if secio {
		log.Printf("Now run \"go run main.go -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fullAddr)
	} else {
		log.Printf("Now run \"go run main.go -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
	}

	//START: writing to shell scripts
	file, err := os.Create("result.sh")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	if secio {
		fmt.Fprintf(file, "#!/bin/sh\n")
		//fmt.Fprintf(file, "xterm -e bash -c 'go run main.go -l %d -d %s -secio'", listenPort+1, fullAddr)
		fmt.Fprintf(file, "go run main.go -l %d -d %s -secio", listenPort+1, fullAddr)
	} else {
		fmt.Fprintf(file, "#!/bin/sh\n")
		//fmt.Fprintf(file, "xterm -e bash -c 'go run main.go -l %d -d %s'", listenPort+1, fullAddr)
		fmt.Fprintf(file, "go run main.go -l %d -d %s", listenPort+1, fullAddr)
	}
	//END: writing to shell scripts

	return basicHost, nil
}

func handleStream(s network.Stream) {

	log.Println("Got a new stream!")

	// Create a buffer stream for non blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

func readData(rw *bufio.ReadWriter) {

	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {

			t1 := tangle{
				Transactions: make([]Transaction, 0),
				Links:        make([]Link, 0),
				State:        make(map[string]int),
			}

			if err := json.Unmarshal([]byte(str), &t1); err != nil {
				log.Fatal(err)
			}

			mutex.Lock()
			if len(t1.Transactions) > len(Tangle.Transactions) &&
				len(t1.Links) > len(Tangle.Links) {
				Tangle.Transactions = t1.Transactions
				Tangle.Links = t1.Links
				Tangle.State = t1.State

				bytes, err := json.MarshalIndent(Tangle, "", "  ")
				if err != nil {

					log.Fatal(err)
				}
				// Green console color: 	\x1b[32m
				// Reset console color: 	\x1b[0m
				fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))

			}
			mutex.Unlock()
		}
	}
}

func writeData(rw *bufio.ReadWriter) {

	go func() {
		for {
			time.Sleep(5 * time.Second)
			mutex.Lock()
			bytes, err := json.Marshal(Tangle)
			if err != nil {
				log.Println(err)
			}
			mutex.Unlock()

			mutex.Lock()
			rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
			rw.Flush()
			mutex.Unlock()

		}
	}()

	stdReader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		sendData = strings.Replace(sendData, "\n", "", -1)

		ss := strings.Fields(sendData)

		amountString := ss[1]
		amountInt, err := strconv.Atoi(amountString)
		if err != nil {
			log.Fatal(err)
		}

		from := ss[3]
		to := ss[5]

		Tangle.State[from] = Tangle.State[from] - amountInt
		Tangle.State[to] = Tangle.State[to] + amountInt

		newTransaction := generateTransaction(Tangle.Transactions[len(Tangle.Transactions)-1], sendData)

		//START: New way of forming Links

		candidates := []int{}
		for _, c := range Tangle.Transactions {
			if newTransaction.TimeInt-Tangle.h > c.TimeInt {
				candidates = append(candidates, c.Index)
			}
		}

		candidateLinks := []Link{}
		for _, l := range Tangle.Links {
			if newTransaction.TimeInt-Tangle.h > Tangle.Transactions[l.Source].TimeInt {
				candidateLinks = append(candidateLinks, l)
			}
		}

		tips := getTips(Tangle.tipSelection, candidates, candidateLinks)
		fmt.Println("NEW TIPS: ", tips)

		mutex.Lock()
		Tangle.Transactions = append(Tangle.Transactions, newTransaction)
		if len(tips) > 0 {
			newLink := generateLink(Tangle.Transactions[tips[0]], newTransaction)
			Tangle.Links = append(Tangle.Links, newLink)
			if len(tips) > 1 && tips[0] != tips[1] {
				newLink := generateLink(Tangle.Transactions[tips[1]], newTransaction)
				Tangle.Links = append(Tangle.Links, newLink)
			}
		}
		mutex.Unlock()

		//END: New way of forming Links

		// newLink1 := generateLink(Tangle.Transactions[len(Tangle.Transactions)-1], newTransaction)
		// newLink2 := generateLink(Tangle.Transactions[len(Tangle.Transactions)-2], newTransaction)

		// if isTransactionValid(newTransaction, Tangle.Transactions[len(Tangle.Transactions)-1]) {
		// 	mutex.Lock()
		// 	Tangle.Transactions = append(Tangle.Transactions, newTransaction)
		// 	Tangle.Links = append(Tangle.Links, newLink1)
		// 	Tangle.Links = append(Tangle.Links, newLink2)
		// 	mutex.Unlock()
		// }

		//START: Good place to print
		calculateWeights()

		//END: Good place to print

		bytes, err := json.Marshal(Tangle)
		if err != nil {
			log.Println(err)
		}

		//spew.Dump(Tangle)

		mutex.Lock()
		rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
		rw.Flush()
		mutex.Unlock()

	}

}

func main() {

	generateTangle()

	// LibP2P code uses golog to log messages. They log with different
	// string IDs (i.e. "swarm"). We can control the verbosity level for
	// all loggers with:
	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info

	// Parse options from the command line
	listenF := flag.Int("l", 0, "wait for incoming connections")
	target := flag.String("d", "", "target peer to dial")
	secio := flag.Bool("secio", false, "enable secio")
	seed := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	if *listenF == 0 {
		log.Fatal("Please provide a port to bind on with -l")
	}

	// Make a host that listens on the given multiaddress
	ha, err := makeBasicHost(*listenF, *secio, *seed)
	if err != nil {
		log.Fatal(err)
	}

	if *target == "" {
		log.Println("listening for connections")
		// Set a stream handler on host A. /p2p/1.0.0 is
		// a user-defined protocol name.
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		select {} // hang forever
		/**** This is where the listener code ends ****/
	} else {
		ha.SetStreamHandler("/p2p/1.0.0", handleStream)

		// The following code extracts target's peer ID from the
		// given multiaddress
		ipfsaddr, err := ma.NewMultiaddr(*target)
		if err != nil {
			log.Fatalln(err)
		}

		pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
		if err != nil {
			log.Fatalln(err)
		}

		peerid, err := peer.IDB58Decode(pid)
		if err != nil {
			log.Fatalln(err)
		}

		// Decapsulate the /ipfs/<peerID> part from the target
		// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
		targetPeerAddr, _ := ma.NewMultiaddr(
			fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
		targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

		// We have a peer ID and a targetAddr so we add it to the peerstore
		// so LibP2P knows how to contact it
		ha.Peerstore().AddAddr(peerid, targetAddr, pstore.PermanentAddrTTL)

		log.Println("opening stream")
		// make a new stream from host B to host A
		// it should be handled on host A by the handler we set above because
		// we use the same /p2p/1.0.0 protocol
		s, err := ha.NewStream(context.Background(), peerid, "/p2p/1.0.0")
		if err != nil {
			log.Fatalln(err)
		}
		// Create a buffered stream so that read and writes are non blocking.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		// Create a thread to read and write data.
		go writeData(rw)
		go readData(rw)

		select {} // hang forever

	}
}

func generateTangle() {
	Tangle.lambda = 1.5
	Tangle.alpha = 0.5
	Tangle.h = 1
	Tangle.tipSelection = "weightedMCMC"

	//Initial State
	state := make(map[string]int)

	state["Alice"] = 50
	state["Bob"] = 50
	state["Charles"] = 50

	Tangle.State = state

	now := time.Now()
	initialTime = now.UnixNano()
	genesisTransaction0 := Transaction{0, "", now.UnixNano() / 1000000, time.Unix(0, now.UnixNano()).String(), 1, 0}
	genesisLink00 := generateLink(genesisTransaction0, genesisTransaction0)

	mutex.Lock()
	Tangle.Transactions = append(Tangle.Transactions, genesisTransaction0)
	Tangle.Links = append(Tangle.Links, genesisLink00)
	mutex.Unlock()

	transactionCount := 10

	now = time.Now()
	myTime := now.UnixNano() / 1000000
	delay := int64(3)

	for len(Tangle.Transactions) < transactionCount {

		myTime = myTime + delay

		newTransaction := Transaction{
			len(Tangle.Transactions),
			"nothing",
			myTime,
			time.Unix(0, myTime*1000000).String(),
			1,
			0}

		mutex.Lock()
		Tangle.Transactions = append(Tangle.Transactions, newTransaction)
		mutex.Unlock()
	}

	for _, t := range Tangle.Transactions {

		candidates := []int{}
		for _, c := range Tangle.Transactions {
			if t.TimeInt-Tangle.h > c.TimeInt {
				candidates = append(candidates, c.Index)
			}
		}

		candidateLinks := []Link{}
		for _, l := range Tangle.Links {
			if t.TimeInt-Tangle.h > Tangle.Transactions[l.Source].TimeInt {
				candidateLinks = append(candidateLinks, l)
			}
		}

		tips := getTips(Tangle.tipSelection, candidates, candidateLinks)
		fmt.Println("START TIPS: ", tips)

		mutex.Lock()
		if len(tips) > 0 {
			newLink := generateLink(Tangle.Transactions[tips[0]], t)
			Tangle.Links = append(Tangle.Links, newLink)
			if len(tips) > 1 && tips[0] != tips[1] {
				newLink := generateLink(Tangle.Transactions[tips[1]], t)
				Tangle.Links = append(Tangle.Links, newLink)
			}
		}
		mutex.Unlock()
	}
}

// make sure transaction is valid by checking index, and comparing the index of the previous transaction
func isTransactionValid(newTransaction, oldTransaction Transaction) bool {
	if oldTransaction.Index+1 != newTransaction.Index {
		return false
	}

	return true
}

// create a new Transaction using previous Transactions index
func generateTransaction(lastTransaction Transaction, Operation string) Transaction {

	var newTransaction Transaction

	newTransaction.Index = lastTransaction.Index + 1
	newTransaction.Operation = Operation

	now := time.Now()
	newTransaction.TimeInt = now.UnixNano() / 1000000
	newTransaction.TimeString = time.Unix(0, now.UnixNano()).String()
	newTransaction.Weight = 1

	return newTransaction
}

func generateLink(target Transaction, source Transaction) Link {

	var newLink Link

	newLink.Target = target.Index
	newLink.Source = source.Index

	return newLink

}

func getApprovedNodes(transaction Transaction) ([]Transaction, []Link) {
	return getDescendants(transaction)
}

func getApprovingNodes(transaction Transaction) ([]Transaction, []Link) {
	return getAncestors(transaction)
}

func isTip(transaction Transaction) bool {

	cuenta := 0

	for _, link := range Tangle.Links {
		if link.Target == transaction.Index {
			cuenta++
		}
	}

	if cuenta < 2 {
		return true
	}

	return false

}

func getTips(algorithm string, candidates []int, candidateLinks []Link) []int {

	if algorithm == "uniformRandom" {

		paso1 := []int{}
		for _, t := range candidates {
			if isTip(Tangle.Transactions[t]) {
				paso1 = append(paso1, t)
			}
		}

		tips := []int{}
		for _, t := range paso1 {
			for _, l := range candidateLinks {
				if l.Source == t {
					tips = append(tips, t)
				}
			}
		}

		if len(tips) == 0 {
			return []int{}
		}
		return []int{choose(tips), choose(tips)}
	}
	if algorithm == "unWeightedMCMC" {

		if len(Tangle.Transactions) == 0 {
			return []int{}
		}

		start := Tangle.Transactions[0]

		return []int{randomWalk(start).Index, randomWalk(start).Index}

	}
	if algorithm == "weightedMCMC" {

		if len(Tangle.Transactions) == 0 {
			return []int{}
		}

		start := Tangle.Transactions[0]

		calculateWeights()

		return []int{weightedRandomWalk(start).Index, weightedRandomWalk(start).Index}

	}
	return []int{}

}

func getAncestors(root Transaction) ([]Transaction, []Link) {

	visitedTransactions := []Transaction{}
	visitedLinks := []Link{}

	stack := []Transaction{root}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		incomingEdges := []Link{}
		for _, link := range Tangle.Links {
			if link.Target == current.Index {
				incomingEdges = append(incomingEdges, link)
			}
		}
		for _, link := range incomingEdges {
			visitedLinks = append(visitedLinks, link)
			yaEsta := false
			for _, transaction := range visitedTransactions {
				if link.Source == transaction.Index {
					yaEsta = true
				}
			}
			if !yaEsta {
				stack = append(stack, Tangle.Transactions[link.Source])
				visitedTransactions = append(visitedTransactions, Tangle.Transactions[link.Source])
			}
		}

	}

	return visitedTransactions, visitedLinks
}

func getDescendants(root Transaction) ([]Transaction, []Link) {

	visitedTransactions := []Transaction{}
	visitedLinks := []Link{}

	stack := []Transaction{root}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		outgoingEdges := []Link{}
		for _, link := range Tangle.Links {
			if link.Source == current.Index {
				outgoingEdges = append(outgoingEdges, link)
			}
		}
		for _, link := range outgoingEdges {
			visitedLinks = append(visitedLinks, link)
			yaEsta := false
			for _, transaction := range visitedTransactions {
				if link.Target == transaction.Index {
					yaEsta = true
				}
			}
			if !yaEsta {
				stack = append(stack, Tangle.Transactions[link.Target])
				visitedTransactions = append(visitedTransactions, Tangle.Transactions[link.Target])
			}
		}

	}

	return visitedTransactions, visitedLinks
}

func choose(array []int) int {
	source := mrand.NewSource(time.Now().UnixNano())
	r := mrand.New(source)
	index := r.Intn(len(array))

	return array[index]
}

func getApprovers(transanction Transaction) []int {
	approvers := []int{}

	for _, link := range Tangle.Links {
		if link.Target == transanction.Index {
			if link.Source != 0 {
				approvers = append(approvers, link.Source)
			}
		}
	}

	return approvers
}

func getChildrenLists() [][]int {

	l := len(Tangle.Transactions)

	childrenLists := make([][]int, l)

	for _, link := range Tangle.Links {
		childrenLists[link.Source] = append(childrenLists[link.Source], link.Target)
	}

	return childrenLists
}

// DFS-based topological sort
func topologicalSort() []int {
	childrenLists := getChildrenLists()
	unvisited := Tangle.Transactions[1:]
	result := []int{}

	for len(unvisited) > 0 {
		t := unvisited[0]
		result, unvisited = visit(t, unvisited, childrenLists, result)
	}

	// Reverse slice
	for i := len(result)/2 - 1; i >= 0; i-- {
		opp := len(result) - 1 - i
		result[i], result[opp] = result[opp], result[i]
	}

	// Add 0
	result = append(result, 0)
	return result
}

func visit(transaction Transaction, unvisited []Transaction, childrenLists [][]int, result []int) ([]int, []Transaction) {

	Esta := false
	for _, t := range unvisited {
		if transaction.Index == t.Index {
			Esta = true
		}
	}
	if !Esta {
		return nil, nil
	}

	for _, child := range childrenLists[transaction.Index] {
		visit(Tangle.Transactions[child], unvisited, childrenLists, result)
	}

	result = append(result, unvisited[0].Index)
	newUnvisited := unvisited[1:]

	return result, newUnvisited
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func randomFloat() float64 {
	source := mrand.NewSource(time.Now().UnixNano())
	r := mrand.New(source)

	return r.Float64()
}

func calculateWeights() {
	sorted := topologicalSort()
	sorted = sorted[:len(sorted)-1]

	//Initialize an empty slice for each node
	l := len(Tangle.Transactions)
	ancestorSlices := make([][]int, l)

	childrenLists := getChildrenLists()

	for _, node := range sorted {
		for _, child := range childrenLists[node] {
			ancestorSlices[child] = append(ancestorSlices[child], ancestorSlices[node]...)
			ancestorSlices[child] = append(ancestorSlices[child], node)
		}
		ancestorSlices[node] = unique(ancestorSlices[node])
		Tangle.Transactions[node].CumWeight = len(ancestorSlices[node]) + 1
	}
}

func unique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func randomWalk(start Transaction) Transaction {
	particle := start

	if !isTip(particle) {
		approvers := getApprovers(particle)
		if len(approvers) != 0 {
			particle = randomWalk(Tangle.Transactions[choose(approvers)])
		}
	}

	return particle
}

func weightedRandomWalk(start Transaction) Transaction {
	particle := start

	if !isTip(particle) {

		approvers := getApprovers(particle)
		if len(approvers) != 0 {

			cumWeights := []int{}
			for _, approver := range approvers {
				cumWeights = append(cumWeights, Tangle.Transactions[approver].CumWeight)
			}

			// normalize so maximum cumWeight is 0
			_, maxCumWeight := minMax(cumWeights)
			normalizedWeights := []int{}
			for i := 0; i < len(cumWeights); i++ {
				normalizedWeights = append(normalizedWeights, cumWeights[i]-maxCumWeight)
			}

			weights := []float64{}
			for i := 0; i < len(normalizedWeights); i++ {
				weights = append(weights, math.Exp(float64(normalizedWeights[i])*float64(Tangle.alpha)))
			}
			myInt := weightedChoose(approvers, weights)
			particle = weightedRandomWalk(Tangle.Transactions[myInt])
		}
	}

	return particle
}

func weightedChoose(approvers []int, weights []float64) int {
	sum := float64(0)
	for i := 0; i < len(weights); i++ {
		sum = sum + weights[i]
	}
	rand := randomFloat() * sum

	cumSum := weights[0]
	for i := 1; i < len(approvers); i++ {
		if rand < cumSum {
			return approvers[i-1]
		}
		cumSum = cumSum + weights[i]
	}
	return approvers[len(approvers)-1]
}

func minMax(array []int) (int, int) {
	var max = array[0]
	var min = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
