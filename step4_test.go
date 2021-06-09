package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testCatalogStep4 = `
{
	"items": [{
			"name": "Apple",
			"price": 0.60,
			"deal": {
				"buy": 1,
				"get": 1
			},
			"stock": 4
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

func TestOrderCostOutOfStock(t *testing.T) {
	s := newTestOrderCartServer(t, []byte(testCatalogStep4))

	obs := testObserver{}
	s.eNotifier.Register(&obs)
	obs.On("OnNotify").Return()

	request := createGetOrderCostRequest([]string{"Apple", "Apple", "Orange", "Apple"})
	_, err := s.GetOrderCost(context.Background(), request)
	assert.Error(t, err)
}
