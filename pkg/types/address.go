package types

import "github.com/btcsuite/btcd/btcutil/base58"

type Address struct {
	value []byte
}

func NewAddressFromByte(value []byte) Address {
	return Address{value}
}

func NewAddressFromBase58String(value string) Address {
	return Address{base58.Decode(value)}
}

func (a Address) AsBytes() []byte {
	return a.value
}

func (a Address) AsBase58String() string {
	return base58.Encode(a.value)
}
