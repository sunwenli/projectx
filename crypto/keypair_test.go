package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignVerifySuccess(t *testing.T) {
	prikey := GeneratePrivateKey()
	pubkey := prikey.PublicKey()

	msg := []byte("hello,world")
	sig, err := prikey.Sign(msg)
	assert.Nil(t, err)

	res := sig.Verify(pubkey, msg)
	assert.True(t, res)
}
func TestSignVerifyFail(t *testing.T) {
	prikey := GeneratePrivateKey()
	pubkey := prikey.PublicKey()

	msg := []byte("hello,world")
	sig, err := prikey.Sign(msg)
	assert.Nil(t, err)

	res := sig.Verify(pubkey, []byte("hhh"))
	assert.False(t, res)

	otherprikey := GeneratePrivateKey()
	otherpubkey := otherprikey.PublicKey()
	res2 := sig.Verify(otherpubkey, msg)
	assert.False(t, res2)
}
