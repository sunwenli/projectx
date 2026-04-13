package core

import (
	"testing"
	"time"

	"github.com/sunwenli/projectx/crypto"
	"github.com/sunwenli/projectx/types"

	"github.com/stretchr/testify/assert"
)

func randomBlock(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privkey := crypto.GeneratePrivateKey()
	tx := randomTransactionWithSignature(t)
	h := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Heigth:        height,
		TimeStamp:     time.Now().UnixNano(),
	}
	b, err := NewBlock(h, []*Transaction{tx})
	assert.Nil(t, err)

	datahash, err := calculateDataHash(b.Transactions)
	assert.Nil(t, err)

	b.Header.DataHash = datahash
	assert.Nil(t, b.Sign(privkey))
	return b
}

func TestBlockSign(t *testing.T) {

	privkey := crypto.GeneratePrivateKey()
	b := randomBlock(t, 0, types.Hash{})
	err := b.Sign(privkey)

	assert.Nil(t, err)
	assert.NotNil(t, b.Signature)

	assert.Nil(t, b.Verify())

	otherprikey := crypto.GeneratePrivateKey()
	b.Validator = otherprikey.PublicKey()

	assert.NotNil(t, b.Verify())

}
func TestVerifyBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(t, 0, types.Hash{})

	assert.Nil(t, b.Sign(privKey))
	assert.Nil(t, b.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t, b.Verify())

	b.Heigth = 100
	assert.NotNil(t, b.Verify())
}
