package types

import (
	"encoding/hex"

	"github.com/cockroachdb/errors"
)

type PubKey struct {
	value []byte
}

func NewPubKeyFromByte(value []byte) PubKey {
	return PubKey{value}
}

func NewPubKeyFromHexString(value string) (PubKey, error) {
	hexVal, err := hex.DecodeString(value)
	if err != nil {
		return PubKey{}, errors.Wrapf(err, "Decoding hex string %v", value)
	}

	return PubKey{hexVal}, nil
}

func (p PubKey) AsBytes() []byte {
	return p.value
}

func (p PubKey) AsHexString() string {
	return hex.EncodeToString(p.value)
}
