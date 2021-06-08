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
	var totalCost float32

	for _, item := range order.GetCart() {
		if itemCost, err := s.inv.GetPrice(item.GetName()); err != nil {
			return nil, err
		} else {
			totalCost += itemCost
		}
	}
	return &pb.OrderCostResponse{Cost: fmt.Sprintf("$%.2f", totalCost)}, nil
}

// Hello returns a greeting for the named person.
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	inv := inventory.Inventory{}
	if err := inv.Load(); err != nil {
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
