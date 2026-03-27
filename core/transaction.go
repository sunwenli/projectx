package core

import (
	"errors"
	"fmt"
	"io"
	"projectx/crypto"
)

type Transaction struct {
	Data      []byte
	From      crypto.PublicKey
	Signature *crypto.Signature
}

func (tx *Transaction) Sign(prikey crypto.PrivateKey) error {
	sig, err := prikey.Sign(tx.Data)
	if err != nil {
		return err
	}

	fmt.Println("sig:", sig)
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

func (t *Transaction) EncodeBinary(w io.Writer) error { return nil }
func (t *Transaction) DecodeBinary(r io.Reader) error { return nil }
