package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"ordercart/inventory"
	"ordercart/notifier"
	"time"

	pb "ordercart/ordercart"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type orderCartServer struct {
	pb.UnimplementedOrderCartServer
	inv       *inventory.Inventory
	eNotifier *notifier.EventNotifier
}

type notifierServer struct {
	pb.UnimplementedCustomerNotificationServer
}

func NewOrderCartServer() (*orderCartServer, error) {
	s := orderCartServer{}
	var err error
	if s.inv, err = inventory.NewInventory(); err != nil {
		return nil, err
	}
	s.eNotifier = notifier.NewEventNotifier()
	return &s, nil
}

func (s *orderCartServer) GetOrderCost(ctx context.Context, order *pb.OrderCostRequest) (*pb.OrderCostResponse, error) {
	id := uuid.New()
	estDeliverTime := time.Now().Add(time.Hour)
	eventRecieved := notifier.NewOrderEvent(id, notifier.Received, estDeliverTime, "")
	s.eNotifier.Notify(eventRecieved)

	cart := make(map[string]uint32)

	for _, item := range order.GetCart() {
		if _, ok := cart[item.GetName()]; ok {
			cart[item.GetName()] += 1
		} else {
			cart[item.GetName()] = 1
		}
	}

	if receipt, totalCost, status, err := s.inv.Process(cart); err != nil {
		s.eNotifier.Notify(notifier.NewOrderEvent(id, status, time.Now(), err.Error()))
		return nil, err
	} else {
		s.eNotifier.Notify(notifier.NewOrderEvent(id, status, estDeliverTime, ""))
		return &pb.OrderCostResponse{Receipt: receipt, Cost: fmt.Sprintf("$%.2f", totalCost)}, nil
	}
}

func (s *notifierServer) OnNotify(e notifier.OrderEvent) {
	switch status := e.Status; status {
	case notifier.FailedInvalidCart, notifier.OutOfStock:
		fmt.Printf("Order id=%s \nstatus=%v\nmessage=%s\n\n", e.OrderID, e.EventStatus(), e.Message)
	case notifier.Received, notifier.Fulfilled:
		fmt.Printf("Order id=%s \nstatus=%v\nest. delivery time=%v\n\n", e.OrderID, e.EventStatus(), e.EstDeliveryTime())
	default:
	}
}

// Hello returns a greeting for the named person.
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	ordCartServer, err := NewOrderCartServer()
	if err != nil {
		log.Fatalf("failed to create new order server, %v", err)
	}
	pb.RegisterOrderCartServer(s, ordCartServer)

	notifyServer := notifierServer{}
	pb.RegisterCustomerNotificationServer(s, &notifyServer)
	ordCartServer.eNotifier.Register(&notifyServer)

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
