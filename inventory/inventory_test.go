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
			},
			"stock": 10
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
	expAppleStock := uint32(10)

	itemPriceTests := []struct {
		name         string
		expPrice     float32
		expDeal      *Offer
		expPriceDeal *PriceOffer
		expStock     *uint32
	}{
		{
			name:     "Apple",
			expPrice: .60,
			expDeal: &Offer{
				Buy: 1,
				Get: 1,
			},
			expStock: &expAppleStock,
		},
		{
			name:     "Orange",
			expPrice: .25,
			expPriceDeal: &PriceOffer{
				Buy:   2,
				Price: 1,
			},
		},
	}

	for _, tt := range itemPriceTests {
		price, err := inv.GetPrice(tt.name)
		assert.NoError(t, err)
		assert.Equal(t, tt.expPrice, price)
		if tt.expDeal != nil {
			deal, err := inv.GetDeal(tt.name)
			assert.NoError(t, err)
			assert.Equal(t, tt.expDeal, deal.(*Offer))
		}
		if tt.expPriceDeal != nil {
			deal, err := inv.GetDeal(tt.name)
			assert.NoError(t, err)
			assert.Equal(t, tt.expPriceDeal, deal.(*PriceOffer))
		}
		if tt.expStock != nil {
			stock, err := inv.GetStock(tt.name)
			assert.NoError(t, err)
			assert.Equal(t, *tt.expStock, *stock)
		}
	}
}
