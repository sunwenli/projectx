package core

import (
	"errors"

	"github.com/sunwenli/projectx/crypto"
	"github.com/sunwenli/projectx/types"
)

type Transaction struct {
	Data      []byte
	From      crypto.PublicKey
	Signature *crypto.Signature

	//cached version of the tx data hash
	hash types.Hash

	//
	firstSeen int64
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}
	return tx.hash
}
func (tx *Transaction) Sign(prikey crypto.PrivateKey) error {
	sig, err := prikey.Sign(tx.Data)
	if err != nil {
		return err
	}

	// fmt.Println("sig:", sig)
	tx.Signature = sig
	tx.From = prikey.PublicKey()
	return nil
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return errors.New("transaction has no signature")
	}
	if !tx.Signature.Verify(tx.From, tx.Data) {
		return errors.New("invalid transaction signature")
	}
	return nil
}

func (tx *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(tx)
}
func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}

func (tx *Transaction) SetFirstSeen(t int64) {
	tx.firstSeen = t
}
func (tx *Transaction) FirstSeen() int64 {
	return tx.firstSeen
}
