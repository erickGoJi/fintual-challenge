package portfolio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenPortfolioRebalanceIsCalledThenItReturnsExpectedOperations(t *testing.T) {
	tests := []struct {
		name           string
		portfolio      Portfolio
		wantOperations []Operation
		wantErr        bool
	}{
		{
			name: "returns buy and sell operations for an unbalanced portfolio",
			portfolio: Portfolio{
				Stocks: []Stock{
					{Symbol: "META", Quantity: 70, Price: 100},
					{Symbol: "AAPL", Quantity: 30, Price: 100},
				},
				Allocations: map[string]float64{"META": 0.40, "AAPL": 0.60},
			},
			wantOperations: []Operation{
				{Symbol: "AAPL", Type: Buy, Quantity: 30, Price: 100, Value: 3000},
				{Symbol: "META", Type: Sell, Quantity: 30, Price: 100, Value: 3000},
			},
		},
		{
			name: "returns no operations for a balanced portfolio",
			portfolio: Portfolio{
				Stocks: []Stock{
					{Symbol: "META", Quantity: 4, Price: 100},
					{Symbol: "AAPL", Quantity: 6, Price: 100},
				},
				Allocations: map[string]float64{"META": 0.4, "AAPL": 0.6},
			},
			wantOperations: []Operation{},
		},
		{
			name: "returns an error for an empty portfolio",
			portfolio: Portfolio{
				Allocations: map[string]float64{"META": 1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			originalStocks := append([]Stock(nil), tt.portfolio.Stocks...)

			// Act
			operations, err := tt.portfolio.Rebalance()

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.wantOperations, operations)
			assert.Equal(t, originalStocks, tt.portfolio.Stocks)
		})
	}
}

func TestPortfolioValidate(t *testing.T) {
	tests := []struct {
		name      string
		portfolio Portfolio
		wantErr   bool
	}{
		{name: "accepts a valid portfolio", portfolio: Portfolio{Stocks: []Stock{{"META", 1, 100}}, Allocations: map[string]float64{"META": 1}}},
		{name: "rejects allocations that do not total one", portfolio: Portfolio{Stocks: []Stock{{"META", 1, 100}}, Allocations: map[string]float64{"META": 0.9}}, wantErr: true},
		{name: "rejects a negative quantity", portfolio: Portfolio{Stocks: []Stock{{"META", -1, 100}}, Allocations: map[string]float64{"META": 1}}, wantErr: true},
		{name: "rejects a zero price", portfolio: Portfolio{Stocks: []Stock{{"META", 1, 0}}, Allocations: map[string]float64{"META": 1}}, wantErr: true},
		{name: "rejects a stock without an allocation", portfolio: Portfolio{Stocks: []Stock{{"META", 1, 100}}, Allocations: map[string]float64{"AAPL": 1}}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			portfolio := tt.portfolio

			// Act
			err := portfolio.Validate()

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestPortfolioTotalValue(t *testing.T) {
	tests := []struct {
		name      string
		portfolio Portfolio
		want      float64
	}{
		{name: "adds all market values", portfolio: Portfolio{Stocks: []Stock{{"META", 2.5, 100}, {"AAPL", 3, 50}}}, want: 400},
		{name: "returns zero for an empty portfolio", portfolio: Portfolio{}, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			portfolio := tt.portfolio

			// Act
			got := portfolio.TotalValue()

			// Assert
			assert.Equal(t, tt.want, got)
		})
	}
}
