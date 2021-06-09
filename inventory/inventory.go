package inventory

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"ordercart/notifier"
)

const catalogFilePath = "ordercart/catalog.json"

type Inventory struct {
	Catalog   map[string]float32
	Deals     map[string]interface{}
	StockPile map[string]*uint32
}

type CatalogView struct {
	Items []Item `json:"items"`
}

type Item struct {
	Name      string      `json:"name"`
	Price     float32     `json:"price"`
	Deal      *Offer      `json:"deal,omitempty"`
	PriceDeal *PriceOffer `json:"priceDeal,omitempty"`
	Stock     *uint32     `json:"stock,omitempty"`
}

type Offer struct {
	Buy uint32 `json:"buy"`
	Get uint32 `json:"get"`
}

type PriceOffer struct {
	Buy   uint32 `json:"buy"`
	Price uint32 `json:"price"`
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
	i.StockPile = make(map[string]*uint32)
	for _, item := range catalog.Items {
		i.Catalog[item.Name] = item.Price
		i.StockPile[item.Name] = item.Stock
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

func (i *Inventory) GetStock(item string) (*uint32, error) {
	if stock, ok := i.StockPile[item]; ok {
		return stock, nil
	}
	return nil, fmt.Errorf("item, %v, does not exist in inventory", item)
}

func (i *Inventory) GetDeal(item string) (interface{}, error) {
	if deal, ok := i.Deals[item]; ok {
		return deal, nil
	}
	return 0, fmt.Errorf("item, %v, does not exist in inventory", item)
}

// Caluculates receipt for cart. If cart contains an item that does not exist
// in catalog, throw invalid cart error.
func (i *Inventory) Process(cart map[string]uint32) (map[string]uint32, float32, notifier.Status, error) {
	var cost float32
	receipt := make(map[string]uint32)

	for item, count := range cart {
		price, err := i.GetPrice(item)
		if err != nil {
			return nil, 0, notifier.FailedInvalidCart, err
		}
		stock, _ := i.GetStock(item)

		switch v := i.Deals[item].(type) {
		case *Offer:
			getCount := (count / v.Buy)
			rem := count % v.Buy
			cost += float32(count) * price
			receipt[item] = getCount*(v.Buy+v.Get) + rem
		case *PriceOffer:
			getCount := (count / v.Buy)
			rem := count % v.Buy
			cost += float32(getCount*v.Price)*price + float32(rem)*price
			receipt[item] = count
		default:
			cost += float32(count) * price
			receipt[item] = count
		}

		if stock != nil && *stock < receipt[item] {
			return nil, 0, notifier.OutOfStock, fmt.Errorf("item, %v, out of stock", item)
		}
	}

	i.updateStockPile(receipt)
	return receipt, cost, notifier.Fulfilled, nil
}

// Update stock inventory in-memory based on provide receipt.
// Ideally, this is a transaction that is committed to a database to have
// data consistency.
func (i *Inventory) updateStockPile(receipt map[string]uint32) {
	for item, count := range receipt {
		stock := i.StockPile[item]
		if stock != nil {
			upStock := (*stock) - count
			i.StockPile[item] = &upStock
		}
	}
}
