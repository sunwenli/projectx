package util

import (
	"math/rand"
	"testing"

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

func NewTransactionWithSignature(t *testing.T, privkey crypto.PrivateKey, size int) *core.Transaction {
	tx := NewRandomTransaction(size)
	err := tx.Sign(privkey)
	assert.Nil(t, err)
	return tx
}
