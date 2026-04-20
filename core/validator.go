package core

import (
	"errors"
	"fmt"
)

var ErrBlockKnown = errors.New("block already known")

type Validator interface {
	ValidatorBlock(*Block) error
}

type BlockValidator struct {
	bc *BlockChain
}

func NewBlockValidator(bc *BlockChain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

func (v *BlockValidator) ValidatorBlock(b *Block) error {
	if v.bc.HasBlock(b.Heigth) {
		return ErrBlockKnown
	}
	if b.Heigth != v.bc.Heigth()+1 {
		return fmt.Errorf("block (%s) with height (%d) is too heigh => current height (%d)", b.Hash(BlockHasher{}), b.Heigth, v.bc.Heigth())
	}
	prevHeader, err := v.bc.GetHeader(b.Heigth - 1)
	if err != nil {
		return err
	}
	hash := BlockHasher{}.Hash(prevHeader)
	if hash != b.PrevBlockHash {
		return fmt.Errorf("the hash of the previous block (%s) is invalid", b.PrevBlockHash)
	}
	if err := b.Verify(); err != nil {
		return err
	}
	return nil
}
