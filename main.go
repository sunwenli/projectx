package main

import (
	"fmt"
	"projectx/network"
	"time"
)

func main() {
	fmt.Println("hello world")
	trlocal := network.NewLocalTransport("LOCAL")
	trremote := network.NewLocalTransport("REMOTE")

	trlocal.Connect(trremote)
	trremote.Connect(trlocal)

	go func() {
		for {
			trremote.SendMessage(trlocal.Addr(), []byte("hello,world"))
			time.Sleep(1 * time.Second)
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{trlocal},
	}
	s := network.NewServer(opts)
	s.Start()
}
