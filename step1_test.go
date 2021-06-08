package main

import (
	"context"
	"ordercart/inventory"
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
func TestGetOrderCostStep1(t *testing.T) {

	inv := inventory.Inventory{}
	assert.NoError(t, inv.Load())
	s := server{inv: inv}

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
		price, err := s.inv.GetPrice(tt.name)
		assert.NoError(t, err)
		assert.Equal(t, tt.expPrice, price)
	}

	request := createGetOrderCostRequest([]string{"Apple", "Apple", "Orange", "Apple"})
	expResponse := pb.OrderCostResponse{Cost: "$2.05"}

	response, err := s.GetOrderCost(context.Background(), request)
	assert.NoError(t, err)
	assert.EqualValues(t, expResponse.Cost, response.Cost)
}

func TestInvalidItemsGetOrderCostStep1(t *testing.T) {
	inv := inventory.Inventory{}
	assert.NoError(t, inv.Load())
	s := server{inv: inv}
	request := createGetOrderCostRequest([]string{"InvalidItemName"})
	_, err := s.GetOrderCost(context.Background(), request)
	assert.Error(t, err)
}
