package util

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/sunwenli/projectx/core"
	"github.com/sunwenli/projectx/crypto"
	"github.com/sunwenli/projectx/types"
)

func RandomBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

func RandomHash() types.Hash {
	return types.HashFromBytes(RandomBytes(32))
}

func NewRandomTransaction(size int) *core.Transaction {
	return core.NewTransaction(RandomBytes(size))
}

func NewRandomTransactionWithSignature(t *testing.T, privkey crypto.PrivateKey, size int) *core.Transaction {
	tx := NewRandomTransaction(size)
	err := tx.Sign(privkey)
	assert.Nil(t, err)
	return tx
}

func NewRandomBlock(t *testing.T, height uint32, prevBlockHash types.Hash) *core.Block {
	txSigner := crypto.GeneratePrivateKey()
	tx := NewRandomTransactionWithSignature(t, txSigner, 100)
	header := &core.Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Heigth:        height,
		TimeStamp:     time.Now().UnixNano(),
	}
	b, err := core.NewBlock(header, []*core.Transaction{tx})
	assert.Nil(t, err)

	datahash, err := core.CalculateDataHash([]*core.Transaction{tx})
	assert.Nil(t, err)

	b.Header.DataHash = datahash
	return b
}

func NewRandomBlockWithSignature(t *testing.T, pk crypto.PrivateKey, height uint32, prevHash types.Hash) *core.Block {
	b := NewRandomBlock(t, height, prevHash)
	assert.Nil(t, b.Sign(pk))

	return b
}
