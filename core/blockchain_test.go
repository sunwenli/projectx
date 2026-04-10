package core

import (
	"fmt"
	"testing"

	"github.com/sunwenli/projectx/types"

	"github.com/stretchr/testify/assert"
)

func newgenesisblock(t *testing.T) *BlockChain {
	bc, err := NewBlockChain(randomBlock(t, 0, types.Hash{}))
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
		assert.Nil(t, bc.AddBlock(randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))))
		// assert.Nil(t, bc.AddBlock(randomBlock(uint32(i+1))))
	}

	assert.Equal(t, len(bc.headers), lenblock+1)
	assert.Equal(t, bc.Heigth(), uint32(lenblock))
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 90, types.Hash{})))
}

func TestGetHeader(t *testing.T) {
	bc := newgenesisblock(t)

	lenblock := 1000
	for i := 0; i < lenblock; i++ {
		b := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(b))
		header, err := bc.GetHeader(b.Heigth)
		assert.Nil(t, err)
		assert.Equal(t, header, b.Header)
	}

	assert.Equal(t, len(bc.headers), lenblock+1)
	assert.Equal(t, bc.Heigth(), uint32(lenblock))
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 90, types.Hash{})))
}

func getPrevBlockHash(t *testing.T, bc *BlockChain, height uint32) types.Hash {
	prevheader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevheader)
}
