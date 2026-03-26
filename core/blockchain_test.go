package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newgenesisblock(t *testing.T) *BlockChain {
	bc, err := NewBlockChain(randomBlock(0))
	assert.Nil(t, err)
	return bc
}
func TestBlockChain(t *testing.T) {
	bc := newgenesisblock(t)
	assert.NotNil(t, bc.validator)
	fmt.Println(bc.Heigth())
}

func TestHasBlock(t *testing.T) {
	bc := newgenesisblock(t)
	assert.True(t, bc.HasBlock(0))
}

func TestAddBlock(t *testing.T) {
	bc := newgenesisblock(t)

	lenblock := 1000
	for i := 0; i < lenblock; i++ {

		assert.Nil(t, bc.AddBlock(randomBlockWithSignature(t, uint32(i+1))))
		// assert.Nil(t, bc.AddBlock(randomBlock(uint32(i+1))))
	}

	assert.Equal(t, len(bc.headers), lenblock+1)
	assert.Equal(t, bc.Heigth(), uint32(lenblock))

	// assert.NotNil(t, bc.AddBlock(randomBlockWithSignature(t, 89)))
	assert.NotNil(t, bc.AddBlock(randomBlock(90)))
}
