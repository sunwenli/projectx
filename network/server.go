package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/sunwenli/projectx/core"
	"github.com/sunwenli/projectx/crypto"
	"github.com/sunwenli/projectx/types"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	Transport     Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
	SeedNodes     []string
	ListenAddr    string
	TCPTransport  *TCPTransport
}
type Server struct {
	TCPTransport *TCPTransport
	peerCh       chan *TCPPeer
	peerMap      map[net.Addr]*TCPPeer
	ServerOpts
	mempool     *TxPool
	chain       *core.BlockChain
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultDecodeRPCFunc
	}
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "addr", opts.ID)
	}
	chain, err := core.NewBlockChain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}
	peerCh := make(chan *TCPPeer)
	tr := NewTCPTransport(opts.ListenAddr, peerCh)

	s := &Server{
		TCPTransport: tr,
		peerCh:       peerCh,
		peerMap:      make(map[net.Addr]*TCPPeer),
		ServerOpts:   opts,
		chain:        chain,
		mempool:      NewTxPool(1000),
		isValidator:  opts.PrivateKey != nil,
		rpcCh:        make(chan RPC),
		quitCh:       make(chan struct{}, 1),
	}
	s.TCPTransport.peerCh = peerCh
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}
	if s.isValidator {
		go s.validatorLoop()
	}

	return s, nil
}

func (s *Server) bootstrapNetwork() {
	for _, addr := range s.SeedNodes {
		fmt.Println("trying to connect to ", addr)

		go func(addr string) {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Printf("could not connect to %+v\n", conn)
				return
			}

			s.peerCh <- &TCPPeer{
				conn: conn,
			}
		}(addr)
	}
}

func (s *Server) Start() {
	s.TCPTransport.Start()
	time.Sleep(time.Second * 1)

	s.bootstrapNetwork()

	s.Logger.Log("msg", "accepting TCP connection on", "addr", s.ListenAddr, "id", s.ID)
free:
	for {
		select {
		case peer := <-s.peerCh:
			// TODO: add mutex PLZ!!!
			s.peerMap[peer.conn.RemoteAddr()] = peer

			go peer.readLoop(s.rpcCh)
			fmt.Printf("new peer => %+v\n", peer)

		case rpc := <-s.rpcCh:
			decodemsg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("error", err)
				continue
			}
			if err := s.RPCProcessor.ProcessMessage(decodemsg); err != nil {
				s.Logger.Log("error", err)
			}
		case <-s.quitCh:
			break free
		}
	}

	fmt.Println("server shutdown")

}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.BlockTime)
	s.Logger.Log("msg", "Starting validator loop", "time", s.BlockTime)
	for {
		<-ticker.C
		s.createNewBlock()
	}
}

func (s *Server) ProcessMessage(msg *DecodeMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	case *core.Block:
		return s.processBlock(t)
	case *GetStatusMessage:
		// return s.processGetStatusMessage(msg.From, t)
	case *StatusMessage:
		// return s.processStatusMessage(msg.From, t)
	case *GetBlockMessage:
		return s.processGetBlockMessage(msg.From, t)
	}
	return nil
}

func (s *Server) processGetBlockMessage(from net.Addr, data *GetBlockMessage) error {
	fmt.Printf("get block message=>%v\n", data)
	return nil
}

// func (s *Server) processStatusMessage(from net.Addr, data *StatusMessage) error {
// 	fmt.Printf("=> received getstatus response msg from %+v => %+v\n", from, data)
// 	if data.CurrentHeigth <= s.chain.Heigth() {
// 		s.Logger.Log("msg", "cannot sync blockheigth to low", "ourheight", s.chain.Heigth(), "theirheight", data.CurrentHeigth, "addr", from)
// 	}
// 	getblockMessage := &GetBlockMessage{
// 		From: s.chain.Heigth(),
// 		To:   0,
// 	}
// 	buf := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buf).Encode(getblockMessage); err != nil {
// 		return err
// 	}
// 	msg := NewMessage(MessageTypeGetBlocks, buf.Bytes())

//		return s.Transport.SendMessage(from, msg.Byte())
//	}
func (s *Server) processGetStatusMessage(from net.Addr, data *GetStatusMessage) error {
	fmt.Printf("=> received getstatus msg from %+v => %+v\n", from, data)

	stausMessage := &StatusMessage{
		ID:            s.ID,
		CurrentHeigth: s.chain.Heigth(),
	}
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(stausMessage); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeStatus, buf.Bytes())

	return s.Transport.SendMessage(from, msg.Byte())
}
func (s *Server) processBlock(b *core.Block) error {
	if err := s.chain.AddBlock(b); err != nil {
		return err
	}

	go s.broadcastBlock(b)
	return nil
}
func (s *Server) processTransaction(tx *core.Transaction) error {
	if err := tx.Verify(); err != nil {
		return err
	}
	hash := tx.Hash(core.TxHasher{})
	if s.mempool.Contains(hash) {
		s.Logger.Log(
			"msg", "transaction already in mempool",
			"hash", hash,
		)
		return nil
	}

	// s.Logger.Log(
	// 	"msg", "adding new tx to mempool",
	// 	"hash", hash,
	// 	"mempoolpending", s.mempool.PendingCount(),
	// )
	go s.broadcasttx(tx)
	s.mempool.Add(tx)
	return nil
}
func (s *Server) broadcast(payload []byte) error {
	for netAddr, peer := range s.peerMap {
		if err := peer.Send(payload); err != nil {
			fmt.Printf("peer send error => addr %s [err: %s]\n", netAddr, err)
		}
	}
	return nil
}

func (s *Server) broadcastBlock(b *core.Block) error {
	buf := &bytes.Buffer{}
	if err := b.Encode(core.NewGobBlockEncoder(buf)); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeBlock, buf.Bytes())

	return s.broadcast(msg.Byte())
}
func (s *Server) broadcasttx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewTxGobEncoder(buf)); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Byte())
}

func (s *Server) createNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Heigth())
	if err != nil {
		return err
	}
	txx := s.mempool.Pending()
	block, err := core.NewBlockFromPrevHeader(currentHeader, txx)
	if err != nil {
		return err
	}
	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}
	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

	s.mempool.PendingClear()

	go s.broadcastBlock(block)
	return nil
}

func genesisBlock() *core.Block {
	header := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Heigth:    0,
		TimeStamp: time.Now().UnixNano(),
	}
	b, _ := core.NewBlock(header, nil)
	return b
}
