package core

import (
	"projectx/crypto"
	"projectx/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func randomBlock(height uint32) *Block {
	h := &Header{
		Version:       1,
		PrevBlockHash: types.RandomHash(),
		Heigth:        height,
		TimeStamp:     time.Now().UnixNano(),
	}
	tx := Transaction{
		Data: []byte("foo"),
	}
	return &Block{
		Header:       h,
		Transactions: []Transaction{tx},
	}
}

func randomBlockWithSignature(t *testing.T, height uint32) *Block {
	prikey := crypto.GeneratePrivateKey()
	b := randomBlock(height)
	err := b.Sign(prikey)
	assert.Nil(t, err)
	return b
}

func TestBlockSign(t *testing.T) {

	privkey := crypto.GeneratePrivateKey()
	b := randomBlock(0)
	err := b.Sign(privkey)

	assert.Nil(t, err)
	assert.NotNil(t, b.Signature)

	assert.Nil(t, b.Verify())

	otherprikey := crypto.GeneratePrivateKey()
	b.Validator = otherprikey.PublicKey()

	assert.NotNil(t, b.Verify())

}
