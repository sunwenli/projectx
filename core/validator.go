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
	if b.Heigth != v.bc.Heigth()+1 {
		return fmt.Errorf("block （%s） too high", b.Hash(BlockHasher{}))
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
