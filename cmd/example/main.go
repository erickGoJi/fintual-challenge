package main

import (
	"fmt"
	"log"

	"github.com/erickGoJi/portfolio-management/portfolio"
)

func main() {
	p := portfolio.Portfolio{
		Stocks: []portfolio.Stock{
			{Symbol: "META", Quantity: 70, Price: 100},
			{Symbol: "AAPL", Quantity: 30, Price: 100},
		},
		Allocations: map[string]float64{"META": 0.40, "AAPL": 0.60},
	}

	operations, err := p.Rebalance()
	if err != nil {
		log.Fatal(err)
	}
	for _, operation := range operations {
		fmt.Printf("%s %.2f shares of %s at $%.2f ($%.2f)\n", operation.Type, operation.Quantity, operation.Symbol, operation.Price, operation.Value)
	}
}
