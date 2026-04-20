package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/sunwenli/projectx/core"
	"github.com/sunwenli/projectx/crypto"

	"github.com/sunwenli/projectx/network"
)

var transports = []network.Transport{
	network.NewLocalTransport("LOCAL"),
}

func main() {
	// trlocal.Connect(trremoteA)
	// trremoteA.Connect(trremoteB)
	// trremoteB.Connect(trremoteC)
	// trremoteB.Connect(trremoteA)
	// trremoteA.Connect(trlocal)

	initRemoteServer(transports)
	localNode := transports[0]
	trLate := network.NewLocalTransport("LATE_NODE")

	// go func() {
	// 	for {
	// 		if err := sendTransaction(trremoteA, trlocal.Addr()); err != nil {
	// 			logrus.Error(err)
	// 		}
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	// if err := sendGetStatusMessage(trremoteA, "REMOTE_B"); err != nil {
	// 	log.Fatal(err)
	// }
	go func() {
		time.Sleep(7 * time.Second)
		lateserver := makeserver(string(trLate.Addr()), trLate, nil)
		go lateserver.Start()
	}()

	privkey := crypto.GeneratePrivateKey()
	localserver := makeserver("LOCAL", localNode, &privkey)
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
		Transport:  tr,
		PrivateKey: pk,
		ID:         id,
		Transports: transports,
	}
	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func sendGetStatusMessage(tr network.Transport, to network.NetAddr) error {
	var (
		getstatusmsg = new(network.GetStatusMessage)
		buf          = new(bytes.Buffer)
	)

	if err := gob.NewEncoder(buf).Encode(getstatusmsg); err != nil {
		return err
	}
	msg := network.NewMessage(network.MessageTypeGetStatus, buf.Bytes())

	return tr.SendMessage(to, msg.Byte())
}
func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privkey := crypto.GeneratePrivateKey()
	// data := []byte{0x01, 0x0a, 0x02, 0x0a, 0x0b}
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f}
	tx := core.NewTransaction(data)
	tx.Sign(privkey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewTxGobEncoder(buf)); err != nil {
		return err
	}
	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())
	return tr.SendMessage(to, msg.Byte())
}
