package core

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"projectx/crypto"
	"projectx/types"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash
	PrevBlockHash types.Hash
	TimeStamp     int64
	Heigth        uint32
	Nonce         uint64
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)
	return buf.Bytes()
}

type Block struct {
	*Header
	Transactions []Transaction

	Validator crypto.PublicKey
	Signature *crypto.Signature
	//cached version of header hash
	hash types.Hash
}

func NewBlock(h *Header, txx []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txx,
	}
}

func (b *Block) Sign(prikey crypto.PrivateKey) error {
	sig, err := prikey.Sign(b.Header.Bytes())
	if err != nil {
		return fmt.Errorf("sign block err:%v", err)
	}
	b.Validator = prikey.PublicKey()
	b.Signature = sig
	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return errors.New("block has no signature")
	}
	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return errors.New("invalid block signature")
	}
	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}
	return nil
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, *tx)
}

func (b *Block) Encode(w io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(w, b)
}

func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(r, b)
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {

	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}

	return b.hash
}
