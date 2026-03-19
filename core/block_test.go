package core

import (
	"bytes"
	"fmt"
	"projectx/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHeader_Encode_Decode(t *testing.T) {

	h := &Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		TimeStamp: int64(time.Now().UnixNano()),
		Heigth:    10,
		Nonce:     989394,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, h.EncodeBinary(buf))

	hDecode := &Header{}

	assert.Nil(t, hDecode.DecodeBinary(buf))

	assert.Equal(t, hDecode, h)
}

func TestBlock_Encode_Decode(t *testing.T) {
	b := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			TimeStamp: int64(time.Now().UnixNano()),
			Heigth:    10,
			Nonce:     989394,
		},
		Transactions: nil,
	}
	buf := &bytes.Buffer{}

	assert.Nil(t, b.EncodeBinary(buf))

	bDecode := &Block{}

	assert.Nil(t, bDecode.DecodeBinary(buf))

	assert.Equal(t, b, bDecode)
}

func TestBlockHash(t *testing.T) {
	b := &Block{
		Header: Header{
			Version:   1,
			PrevBlock: types.RandomHash(),
			TimeStamp: int64(time.Now().UnixNano()),
			Heigth:    10,
			Nonce:     989394,
		},
		Transactions: []Transaction{},
	}
	h := b.Hash()
	fmt.Println(h.String())
	assert.False(t, h.IsZero())
}
