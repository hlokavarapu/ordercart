package inventory

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const catalogFilePath = "ordercart/catalog.json"

type Inventory struct {
	Catalog map[string]float32
	Deals   map[string]interface{}
}

type CatalogView struct {
	Items []Item `json:"items"`
}

type Item struct {
	Name      string      `json:"name"`
	Price     float32     `json:"price"`
	Deal      *Offer      `json:"deal,omitempty"`
	PriceDeal *PriceOffer `json:"priceDeal,omitempty"`
}

type Offer struct {
	Buy int32 `json:"buy"`
	Get int32 `json:"get"`
}

type PriceOffer struct {
	Buy   int32 `json:"buy"`
	Price int32 `json:"price"`
}

func NewInventory() (*Inventory, error) {
	inv := Inventory{}
	catalogData, err := ioutil.ReadFile(catalogFilePath)
	if err != nil {
		return nil, err
	}
	inv.Load(catalogData)
	return &inv, nil
}

func (i *Inventory) Load(catalogData []byte) error {
	catalog := CatalogView{}
	if err := json.Unmarshal(catalogData, &catalog); err != nil {
		return err
	}

	i.Catalog = make(map[string]float32)
	i.Deals = make(map[string]interface{})
	for _, item := range catalog.Items {
		i.Catalog[item.Name] = item.Price
		if item.Deal != nil {
			i.Deals[item.Name] = item.Deal
		} else if item.PriceDeal != nil {
			i.Deals[item.Name] = item.PriceDeal
		}
	}

	return nil
}

// Get item price from catalog. If item does not exist in catalog, throw error.
func (i *Inventory) GetPrice(item string) (float32, error) {
	if cost, ok := i.Catalog[item]; ok {
		return cost, nil
	}
	return 0, fmt.Errorf("item, %v, does not exist in inventory", item)
}

// Caluculates receipt for cart. If cart contains an item that does not exist
// in catalog, throw invalid cart error.
func (i *Inventory) GetCost(cart map[string]uint32) (map[string]uint32, float32, error) {
	var cost float32
	receipt := make(map[string]uint32)

	for item, count := range cart {
		price, err := i.GetPrice(item)
		if err != nil {
			return nil, 0, err
		}

		switch v := i.Deals[item].(type) {
		case *Offer:
			getCount := (count / uint32(v.Buy))
			rem := count % uint32(v.Buy)
			cost += float32(count) * price
			receipt[item] = getCount*(uint32(v.Buy)+uint32(v.Get)) + rem
		case *PriceOffer:
			getCount := (count / uint32(v.Buy))
			rem := count % uint32(v.Buy)
			cost += float32(getCount*uint32(v.Price))*price + float32(rem)*price
			receipt[item] = count
		default:
			cost += float32(count) * price
			receipt[item] = count
		}
	}
	return receipt, cost, nil
}
