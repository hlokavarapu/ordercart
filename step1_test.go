package main

import (
	"context"
	"testing"

	pb "ordercart/ordercart"

	"github.com/stretchr/testify/assert"
)

func createGetOrderCostRequest(names []string) *pb.OrderCostRequest {
	cart := []*pb.Item{}
	for _, name := range names {
		cart = append(cart, &pb.Item{Name: name})
	}
	return &pb.OrderCostRequest{Cart: cart}
}

const testCatalogStep1 = `
{
	"items": [{
			"name": "Apple",
			"price": 0.60
		},
		{
			"name": "Orange",
			"price": 0.25
		}
	]
}`

func newTestOrderCartServer(t *testing.T, catalogData []byte) *orderCartServer {
	s, err := NewOrderCartServer()
	assert.NoError(t, err)
	assert.NoError(t, s.inv.Load(catalogData))
	return s
}

func TestGetOrderCost_Step1(t *testing.T) {
	s := newTestOrderCartServer(t, []byte(testCatalogStep1))

	request := createGetOrderCostRequest([]string{"Apple", "Apple", "Orange", "Apple"})
	expResponse := pb.OrderCostResponse{
		Receipt: map[string]uint32{
			"Apple":  3,
			"Orange": 1,
		},
		Cost: "$2.05",
	}

	response, err := s.GetOrderCost(context.Background(), request)
	assert.NoError(t, err)
	assert.EqualValues(t, expResponse.Receipt, response.Receipt)
	assert.EqualValues(t, expResponse.Cost, response.Cost)
}

func TestInvalidItemsGetOrderCostStep1(t *testing.T) {
	s := newTestOrderCartServer(t, []byte(testCatalogStep1))

	request := createGetOrderCostRequest([]string{"InvalidItemName"})
	_, err := s.GetOrderCost(context.Background(), request)
	assert.Error(t, err)
}
