package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testCatalogStep2 = `
{
	"items": [{
			"name": "Apple",
			"price": 0.60,
			"deal": {
				"buy": 1,
				"get": 1
			}
		},
		{
			"name": "Orange",
			"price": 0.25,
			"priceDeal": {
				"buy": 3,
				"price": 2
			}
		}
	]
}`

func TestGetOrderCost_Step2(t *testing.T) {
	s := newTestOrderCartServer(t, []byte(testCatalogStep2))

	cartTests := []struct {
		cart       []string
		expReceipt map[string]uint32
		expCost    string
	}{
		{
			cart: []string{"Apple", "Apple", "Orange", "Apple"},
			expReceipt: map[string]uint32{
				"Apple":  6,
				"Orange": 1,
			},
			expCost: "$2.05",
		},
		{
			cart: []string{"Orange", "Orange", "Orange"},
			expReceipt: map[string]uint32{
				"Orange": 3,
			},
			expCost: "$0.50",
		},
		{
			cart: []string{"Apple", "Apple", "Apple"},
			expReceipt: map[string]uint32{
				"Apple": 6,
			},
			expCost: "$1.80",
		},
	}

	for _, tt := range cartTests {
		request := createGetOrderCostRequest(tt.cart)
		response, err := s.GetOrderCost(context.Background(), request)
		assert.NoError(t, err)
		assert.EqualValues(t, tt.expReceipt, response.Receipt)
		assert.EqualValues(t, tt.expCost, response.Cost)
	}
}
