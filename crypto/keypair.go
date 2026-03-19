package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math/big"
	"projectx/types"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

func NewPrivateKeyFromReader(r io.Reader) PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), r)
	if err != nil {
		panic(err)
	}
	return PrivateKey{
		key: key,
	}
}

func GeneratePrivateKey() PrivateKey {
	return NewPrivateKeyFromReader(rand.Reader)
}

type PublicKey []byte

func (p *PrivateKey) PublicKey() PublicKey {
	return elliptic.MarshalCompressed(p.key.PublicKey, p.key.X, p.key.Y)
}

func (p PublicKey) String() string {
	return hex.EncodeToString(p)
}

func (p PublicKey) Address() types.Address {
	h := sha256.Sum256(p)
	return types.AddressFromByte(h[len(h)-20:])
}

func (p PrivateKey) Sign(data []byte) (*signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, p.key, data)
	if err != nil {
		return nil, err
	}
	sig := &signature{
		R: r,
		S: s,
	}
	return sig, nil
}

type signature struct {
	R, S *big.Int
}

func (sig signature) String() string {
	b := append(sig.S.Bytes(), sig.R.Bytes()...)
	return hex.EncodeToString(b)
}

func (sig signature) Verify(pubkey PublicKey, data []byte) bool {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), pubkey)
	key := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	return ecdsa.Verify(key, data, sig.R, sig.S)
}
