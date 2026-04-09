package network

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sunwenli/projectx/core"
	"github.com/sunwenli/projectx/crypto"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	RPCHandler RPCHandler
	Transports []Transport
	BlockTime  time.Duration
	PrivateKey *crypto.PrivateKey
}
type Server struct {
	ServerOpts
	blockTime   time.Duration
	memPool     *TxPool
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	s := &Server{
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}

	if opts.RPCHandler == nil {
		opts.RPCHandler = NewDefaultRPCHandler(s)
	}
	s.ServerOpts = opts
	return s
}

func (s *Server) Start() {
	s.initTransport()
	ticker := time.NewTicker(s.blockTime)
free:
	for {
		select {
		case rpc := <-s.rpcCh:
			if err := s.RPCHandler.HandleRPC(rpc); err != nil {
				logrus.Error(err)
			}
		case <-s.quitCh:
			break free
		case <-ticker.C:
			if s.isValidator {
				s.createNewBlock()
			}
		}
	}

	fmt.Println("server shutdown")

}

func (s *Server) ProcessTransaction(from NetAddr, tx *core.Transaction) error {
	if err := tx.Verify(); err != nil {
		return err
	}
	hash := tx.Hash(core.TxHasher{})
	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash": hash,
		}).Info("transaction already in mempool")
		return nil
	}
	tx.SetFirstSeen(time.Now().UnixNano())
	logrus.WithFields(logrus.Fields{
		"hash":        hash,
		"mempool len": s.memPool.Len(),
	}).Info("adding new tx to mempool")

	//todo broadcast tx to peers
	return s.memPool.Add(tx)
}
func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
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
