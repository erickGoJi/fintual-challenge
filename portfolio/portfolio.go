package portfolio

type Portfolio struct {
	Stocks      []Stock
	Allocations map[string]float64
}

// TotalValue returns the current market value of the portfolio.
func (p Portfolio) TotalValue() float64 {
	var total float64
	for _, stock := range p.Stocks {
		total += stock.MarketValue()
	}
	return total
}
