package inventory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testCatalog = `
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
				"buy": 2,
				"price": 1
			}
		}
	]
}`

func TestInventoryLoad(t *testing.T) {
	inv := Inventory{}
	assert.NoError(t, inv.Load([]byte(testCatalog)))

	itemPriceTests := []struct {
		name     string
		expPrice float32
	}{
		{
			name:     "Apple",
			expPrice: .60,
		},
		{
			name:     "Orange",
			expPrice: .25,
		},
	}

	for _, tt := range itemPriceTests {
		price, err := inv.GetPrice(tt.name)
		assert.NoError(t, err)
		assert.Equal(t, tt.expPrice, price)
	}
}
