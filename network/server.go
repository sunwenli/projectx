package network

import (
	"bytes"
	"fmt"
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
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}
type Server struct {
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
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}
	chain, err := core.NewBlockChain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}
	s := &Server{
		ServerOpts:  opts,
		chain:       chain,
		mempool:     NewTxPool(1000),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}

	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}
	if s.isValidator {
		go s.validatorLoop()
	}
	return s, nil
}

func (s *Server) Start() {
	s.initTransport()
free:
	for {
		select {
		case rpc := <-s.rpcCh:
			decodemsg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("error", err)
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
	}
	return nil
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
	for _, tr := range s.Transports {
		if err := tr.BroadCast(payload); err != nil {
			return err
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

func (s *Server) initTransport() {

	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)
	}
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
