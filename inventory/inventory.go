package inventory

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const catalogFilePath = "ordercart/catalog.json"

type Inventory struct {
	Catalog map[string]float32
}

type CatalogView struct {
	Items []Item `json:"items"`
}

type Item struct {
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

func (i *Inventory) Load() error {
	catalogData, err := ioutil.ReadFile(catalogFilePath)
	if err != nil {
		return err
	}

	catalog := CatalogView{}
	if err = json.Unmarshal(catalogData, &catalog); err != nil {
		return err
	}

	i.Catalog = make(map[string]float32)
	for _, item := range catalog.Items {
		i.Catalog[item.Name] = item.Price
	}

	return nil
}

func (i *Inventory) GetPrice(item string) (float32, error) {
	if cost, ok := i.Catalog[item]; ok {
		return cost, nil
	}
	return 0, fmt.Errorf("item, %v, does not exist in inventory", item)
}
