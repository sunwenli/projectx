package core

import (
	"fmt"
	"sync"

	"github.com/go-kit/log"
)

type BlockChain struct {
	logger    log.Logger
	headers   []*Header
	lock      sync.RWMutex
	store     Storage
	validator Validator
}

func NewBlockChain(l log.Logger, genesis *Block) (*BlockChain, error) {
	bc := &BlockChain{
		headers: []*Header{},
		store:   NewMemorystore(),
		logger:  l,
	}
	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)
	return bc, err
}

func (bc *BlockChain) SetValidator(v Validator) {
	bc.validator = v
}
func (bc *BlockChain) AddBlock(b *Block) error {
	if err := bc.validator.ValidatorBlock(b); err != nil {
		return err
	}
	for _, tx := range b.Transactions {
		bc.logger.Log("msg", "excuting code", "len", len(tx.Data), "hash", tx.Hash(&TxHasher{}))

		vm := NewVm(tx.Data)
		if err := vm.Run(); err != nil {
			return err
		}

		bc.logger.Log("vm result", vm.stack[vm.sp])
	}
	return bc.addBlockWithoutValidation(b)
}

func (bc *BlockChain) HasBlock(heigth uint32) bool {
	return heigth <= bc.Heigth()
}

func (bc *BlockChain) GetHeader(heigth uint32) (*Header, error) {
	if heigth > bc.Heigth() {
		return nil, fmt.Errorf("given height (%d) too high", heigth)
	}
	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[heigth], nil
}
func (bc *BlockChain) Heigth() uint32 {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	return uint32(len(bc.headers) - 1)
}

func (bc *BlockChain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.lock.Unlock()

	bc.logger.Log(
		"msg", "adding new block",
		"height", b.Heigth,
		"hash", b.Hash(BlockHasher{}),
	)
	return bc.store.Put(b)
}
