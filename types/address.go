package types

import (
	"encoding/hex"
	"fmt"
)

type Address [20]uint8

func (a Address) ToSlice() []byte {
	value := make([]byte, 20)
	for i := 0; i < 20; i++ {
		value[i] = a[i]
	}
	return value
}

func (a Address) String() string {
	return hex.EncodeToString(a.ToSlice())
}

func AddressFromByte(b []byte) Address {
	if len(b) != 20 {
		msg := fmt.Sprintf("given bytes with length %d should be 20", len(b))
		panic(msg)
	}
	address := make([]byte, 20)

	for i := 0; i < 20; i++ {
		address[i] = b[i]
	}

	return Address(address)
}
