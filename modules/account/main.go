package account

type Balance struct {
	Value float64
}

func NewBalance() *Balance {
	return &Balance{Value:1000}
}

func (b *Balance) Inc(value int) *Balance {
	b.Value = b.Value + float64(value)
	return b
}

func (b *Balance) Dec(value int) *Balance {
	if b.Value < float64(value) {
		b.Value = 0
	}
	b.Value = b.Value - float64(value)
	return b
}

func (b *Balance) Current() float64 {
	return b.Value
}
