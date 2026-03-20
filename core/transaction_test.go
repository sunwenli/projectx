package core

import (
	"projectx/crypto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionSign(t *testing.T) {
	prikey := crypto.GeneratePrivateKey()

	tx := &Transaction{
		Data: []byte("foo"),
	}
	err := tx.Sign(&prikey)
	assert.Nil(t, err)
	assert.NotNil(t, tx.Signature)

	assert.Nil(t, tx.Verify())
	otherprikey := crypto.GeneratePrivateKey()
	tx.Publickey = otherprikey.PublicKey()
	assert.NotNil(t, tx.Verify())
}
