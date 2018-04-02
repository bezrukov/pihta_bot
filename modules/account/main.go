package account

type Balance struct {
	Value float64
}

func (b *Balance) addBalance(value int) *Balance {
	b.Value = b.Value + float64(value)
	return b
}

