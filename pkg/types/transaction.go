package types

import (
	"encoding/hex"

	"github.com/cockroachdb/errors"
)

type Transaction struct {
	Id []byte
}

func NewTransactionFromByte(id []byte) Transaction {
	return Transaction{id}
}

func NewTransactionFromHexString(id string) (Transaction, error) {
	hexVal, err := hex.DecodeString(id)
	if err != nil {
		return Transaction{}, errors.Wrapf(err, "Decoding hex string %v", id)
	}
	return Transaction{hexVal}, nil
}

func (t Transaction) IdAsBytes() []byte {
	return t.Id
}

func (t Transaction) IdAsHexString() string {
	return hex.EncodeToString(t.Id)
}
