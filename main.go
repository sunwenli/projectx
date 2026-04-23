package main

import (
	"bytes"
	"log"
	"net"

	"github.com/sunwenli/projectx/core"
	"github.com/sunwenli/projectx/crypto"

	"github.com/sunwenli/projectx/network"
)

func main() {
	privkey := crypto.GeneratePrivateKey()
	localNode := makeserver("LOCAL_NODE", &privkey, ":3000", []string{":4000"})
	go localNode.Start()

	remoteNode := makeserver("REMOTE_NODE", nil, ":4000", []string{":5000"})
	go remoteNode.Start()

	remoteNodeB := makeserver("REMOTE_NODE_B", nil, ":5000", nil)
	go remoteNodeB.Start()

	tcpTester()
	select {}
}
func makeserver(id string, pk *crypto.PrivateKey, addr string, seedNodes []string) *network.Server {

	opts := network.ServerOpts{
		SeedNodes:  seedNodes,
		ListenAddr: addr,
		PrivateKey: pk,
		ID:         id,
	}
	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func tcpTester() {
	conn, err := net.Dial("tcp", ":3000")

	if err != nil {
		panic(err)
	}
	privKey := crypto.GeneratePrivateKey()
	// data := []byte{0x03, 0x0a, 0x02, 0x0a, 0x0e}
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f}
	tx := core.NewTransaction(data)
	tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewTxGobEncoder(buf)); err != nil {
		panic(err)
	}
	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())
	_, err = conn.Write(msg.Byte())
	if err != nil {
		panic(err)
	}
}
