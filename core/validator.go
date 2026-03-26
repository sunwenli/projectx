package core

import (
	"fmt"
)

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
		return fmt.Errorf("blockchain already has block (%d) with hash (%s)", b.Heigth, b.Hash(BlockHasher{}))
	}
	if err := b.Verify(); err != nil {
		return err
	}
	return nil
}
