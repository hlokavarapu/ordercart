package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"ordercart/inventory"

	pb "ordercart/ordercart"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedOrderCartServer
	inv inventory.Inventory
}

func (s *server) GetOrderCost(ctx context.Context, order *pb.OrderCostRequest) (*pb.OrderCostResponse, error) {
	cart := make(map[string]uint32)

	for _, item := range order.GetCart() {
		if _, ok := cart[item.GetName()]; ok {
			cart[item.GetName()] += 1
		} else {
			cart[item.GetName()] = 1
		}
	}

	receipt, totalCost, err := s.inv.GetCost(cart)
	if err != nil {
		return nil, err
	}
	return &pb.OrderCostResponse{Receipt: receipt, Cost: fmt.Sprintf("$%.2f", totalCost)}, nil
}

// Hello returns a greeting for the named person.
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	inv := inventory.Inventory{}
	if err := inv.LoadCatalog(); err != nil {
		log.Fatalf("failed to get catalog, %v", err)
	}

	ordCartServer := server{
		inv: inv,
	}

	pb.RegisterOrderCartServer(s, &ordCartServer)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
