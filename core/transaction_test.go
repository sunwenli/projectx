package core

import (
	"bytes"
	"testing"

	"github.com/sunwenli/projectx/crypto"

	"github.com/stretchr/testify/assert"
)

func TestTransactionSign(t *testing.T) {
	prikey := crypto.GeneratePrivateKey()

	tx := &Transaction{
		Data: []byte("foo"),
	}
	err := tx.Sign(prikey)
	assert.Nil(t, err)
	assert.NotNil(t, tx.Signature)

	assert.Nil(t, tx.Verify())
	otherprikey := crypto.GeneratePrivateKey()
	tx.From = otherprikey.PublicKey()
	assert.NotNil(t, tx.Verify())
}

func randomTransactionWithSignature(t *testing.T) Transaction {
	prikey := crypto.GeneratePrivateKey()
	tx := Transaction{
		Data: []byte("foo"),
	}
	assert.Nil(t, tx.Sign(prikey))
	return tx
}
func TestR(t *testing.T) {
	randomTransactionWithSignature(t)
}

func TestEncoderDecoder(t *testing.T) {
	tx := randomTransactionWithSignature(t)
	buf := &bytes.Buffer{}

	assert.Nil(t, tx.Encode(NewTxGobEncoder(buf)))

	txdec := new(Transaction)
	assert.Nil(t, txdec.Decode(NewTxGobDecoder(buf)))
	assert.Equal(t, &tx, txdec)
}
