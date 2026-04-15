package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sunwenli/projectx/core"
	"github.com/sunwenli/projectx/crypto"

	"github.com/sunwenli/projectx/network"
)

func main() {
	fmt.Println("hello world")
	trlocal := network.NewLocalTransport("LOCAL")
	trremoteA := network.NewLocalTransport("REMOTE_A")
	trremoteB := network.NewLocalTransport("REMOTE_B")
	trremoteC := network.NewLocalTransport("REMOTE_C")

	trlocal.Connect(trremoteA)
	trremoteA.Connect(trremoteB)
	trremoteB.Connect(trremoteC)
	trremoteA.Connect(trlocal)

	initRemoteServer([]network.Transport{trremoteA, trremoteB, trremoteC})
	go func() {
		for {
			if err := sendTransaction(trremoteA, trlocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		time.Sleep(7 * time.Second)
		trlate := network.NewLocalTransport("LATE_REMOTE")
		trremoteC.Connect(trlate)

		lateserver := makeserver(string(trlate.Addr()), trlate, nil)
		go lateserver.Start()
	}()

	privkey := crypto.GeneratePrivateKey()
	localserver := makeserver("LOCAL", trlocal, &privkey)
	localserver.Start()
}
func initRemoteServer(trs []network.Transport) {
	for i := 0; i < len(trs); i++ {
		id := fmt.Sprintf("REMOTE_%d", i)
		s := makeserver(id, trs[i], nil)
		go s.Start()
	}
}
func makeserver(id string, tr network.Transport, pk *crypto.PrivateKey) *network.Server {

	opts := network.ServerOpts{
		PrivateKey: pk,
		ID:         id,
		Transports: []network.Transport{tr},
	}
	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privkey := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(int64(rand.Intn(100000)), 10))
	tx := core.NewTransaction(data)
	tx.Sign(privkey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewTxGobEncoder(buf)); err != nil {
		return err
	}
	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())
	return tr.SendMessage(to, msg.Byte())
}
