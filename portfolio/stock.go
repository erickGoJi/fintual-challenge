package portfolio

type Stock struct {
	Symbol   string
	Quantity float64
	Price    float64
}

// CurrentPrice returns the last available market price.
func (s Stock) CurrentPrice() float64 {
	return s.Price
}

// MarketValue calculates the total value of a position.
func (s Stock) MarketValue() float64 {
	return s.Quantity * s.CurrentPrice()
}
