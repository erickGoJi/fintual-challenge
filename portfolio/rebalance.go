package portfolio

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
)

const allocationTolerance = 1e-9

type OperationType string

const (
	Buy  OperationType = "BUY"
	Sell OperationType = "SELL"
)

// Operation describes a proposed trade. Rebalance never executes it.
type Operation struct {
	Symbol   string
	Type     OperationType
	Quantity float64
	Price    float64
	Value    float64
}

// Rebalance calculates the trades needed to reach the target allocations.
// Fractional shares are allowed and the receiver is not modified.
func (p Portfolio) Rebalance() ([]Operation, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}

	totalValue := p.TotalValue()
	if totalValue == 0 {
		return []Operation{}, nil
	}

	stocks := make(map[string]Stock, len(p.Stocks))
	for _, stock := range p.Stocks {
		stocks[stock.Symbol] = stock
	}

	symbols := make([]string, 0, len(p.Allocations))
	for symbol := range p.Allocations {
		symbols = append(symbols, symbol)
	}
	sort.Strings(symbols)

	operations := make([]Operation, 0, len(symbols))
	for _, symbol := range symbols {
		stock := stocks[symbol]
		difference := totalValue*p.Allocations[symbol] - stock.MarketValue()
		if math.Abs(difference) <= allocationTolerance {
			continue
		}

		operationType := Buy
		if difference < 0 {
			operationType = Sell
		}
		value := math.Abs(difference)
		operations = append(operations, Operation{
			Symbol:   symbol,
			Type:     operationType,
			Quantity: value / stock.CurrentPrice(),
			Price:    stock.CurrentPrice(),
			Value:    value,
		})
	}

	return operations, nil
}

// Validate checks that the portfolio can be rebalanced safely.
func (p Portfolio) Validate() error {
	if len(p.Allocations) == 0 {
		return errors.New("allocations are required")
	}

	allocationTotal := 0.0
	for symbol, allocation := range p.Allocations {
		if strings.TrimSpace(symbol) == "" {
			return errors.New("allocation symbol is required")
		}
		if math.IsNaN(allocation) || math.IsInf(allocation, 0) || allocation < 0 || allocation > 1 {
			return fmt.Errorf("allocation for %s must be between 0 and 1", symbol)
		}
		allocationTotal += allocation
	}
	if math.Abs(allocationTotal-1) > allocationTolerance {
		return fmt.Errorf("allocations must add up to 1, got %.10g", allocationTotal)
	}

	seen := make(map[string]struct{}, len(p.Stocks))
	for _, stock := range p.Stocks {
		if strings.TrimSpace(stock.Symbol) == "" {
			return errors.New("stock symbol is required")
		}
		if _, exists := seen[stock.Symbol]; exists {
			return fmt.Errorf("stock %s is duplicated", stock.Symbol)
		}
		seen[stock.Symbol] = struct{}{}
		if _, exists := p.Allocations[stock.Symbol]; !exists {
			return fmt.Errorf("stock %s has no target allocation", stock.Symbol)
		}
		if math.IsNaN(stock.Quantity) || math.IsInf(stock.Quantity, 0) || stock.Quantity < 0 {
			return fmt.Errorf("quantity for %s must be non-negative", stock.Symbol)
		}
		if math.IsNaN(stock.Price) || math.IsInf(stock.Price, 0) || stock.Price <= 0 {
			return fmt.Errorf("price for %s must be greater than zero", stock.Symbol)
		}
	}

	for symbol := range p.Allocations {
		if _, exists := seen[symbol]; !exists {
			return fmt.Errorf("stock %s is missing a current price", symbol)
		}
	}
	return nil
}
