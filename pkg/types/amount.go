package types

type Amount struct {
	amount uint64 // in millisats
}

func NewAmountSats(sats uint64) Amount {
	return Amount{sats * 1000}
}

func NewBAmountMSats(msats uint64) Amount {
	return Amount{msats}
}

func (a Amount) AsSats() uint64 {
	return a.amount / 1000 // any millisats left over are discarded
}

func (a Amount) AsMSats() uint64 {
	return a.amount
}
