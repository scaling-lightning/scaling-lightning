package network

type BitcoinAmount struct {
	amount uint64 // in millisats
}

func NewBitcoinAmountSats(sats uint64) BitcoinAmount {
	return BitcoinAmount{sats * 1000}
}

func NewBitcoinAmountMSats(msats uint64) BitcoinAmount {
	return BitcoinAmount{msats}
}

func (a BitcoinAmount) ToSats() uint64 {
	return a.amount / 1000
}

func (a BitcoinAmount) ToMSats() uint64 {
	return a.amount
}
