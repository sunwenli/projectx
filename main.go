package main

import (
	"bytes"
	"fmt"
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
	trremote := network.NewLocalTransport("REMOTE")

	trlocal.Connect(trremote)
	trremote.Connect(trlocal)

	go func() {
		for {
			if err := sendTransaction(trremote, trlocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{trlocal},
	}
	s := network.NewServer(opts)
	s.Start()
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
