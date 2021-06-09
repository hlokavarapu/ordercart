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
	event := notifier.OrderEvent{
		OrderID:      id,
		Status:       notifier.Received,
		Deliverytime: time.Now().Add(time.Hour),
	}
	s.eNotifier.Notify(event)

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
		event.Status = notifier.FailedIvalidCart
		s.eNotifier.Notify(event)
		return nil, err
	} else {
		event.Status = notifier.Fulfilled
		event.Deliverytime = time.Now()
		s.eNotifier.Notify(event)
	}
	return &pb.OrderCostResponse{Receipt: receipt, Cost: fmt.Sprintf("$%.2f", totalCost)}, nil
}

func (s *notifierServer) OnNotify(e notifier.OrderEvent) {
	switch status := e.Status; status {
	case notifier.FailedIvalidCart:
		fmt.Printf("*** Order id=%s \nstatus=%v", e.OrderID, e.EventStatus())
	case notifier.Received | notifier.Fulfilled:
		fmt.Printf("*** Order id=%s \nstatus=%v\nest. delivery time=%v", e.OrderID, e.EventStatus(), e.EstDeliveryTime())
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
