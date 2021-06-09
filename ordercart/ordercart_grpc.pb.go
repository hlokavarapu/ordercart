// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package ordercart

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// OrderCartClient is the client API for OrderCart service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OrderCartClient interface {
	GetOrderCost(ctx context.Context, in *OrderCostRequest, opts ...grpc.CallOption) (*OrderCostResponse, error)
}

type orderCartClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderCartClient(cc grpc.ClientConnInterface) OrderCartClient {
	return &orderCartClient{cc}
}

func (c *orderCartClient) GetOrderCost(ctx context.Context, in *OrderCostRequest, opts ...grpc.CallOption) (*OrderCostResponse, error) {
	out := new(OrderCostResponse)
	err := c.cc.Invoke(ctx, "/ordercart.OrderCart/GetOrderCost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrderCartServer is the server API for OrderCart service.
// All implementations must embed UnimplementedOrderCartServer
// for forward compatibility
type OrderCartServer interface {
	GetOrderCost(context.Context, *OrderCostRequest) (*OrderCostResponse, error)
	mustEmbedUnimplementedOrderCartServer()
}

// UnimplementedOrderCartServer must be embedded to have forward compatible implementations.
type UnimplementedOrderCartServer struct {
}

func (UnimplementedOrderCartServer) GetOrderCost(context.Context, *OrderCostRequest) (*OrderCostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrderCost not implemented")
}
func (UnimplementedOrderCartServer) mustEmbedUnimplementedOrderCartServer() {}

// UnsafeOrderCartServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OrderCartServer will
// result in compilation errors.
type UnsafeOrderCartServer interface {
	mustEmbedUnimplementedOrderCartServer()
}

func RegisterOrderCartServer(s grpc.ServiceRegistrar, srv OrderCartServer) {
	s.RegisterService(&OrderCart_ServiceDesc, srv)
}

func _OrderCart_GetOrderCost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrderCostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderCartServer).GetOrderCost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ordercart.OrderCart/GetOrderCost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderCartServer).GetOrderCost(ctx, req.(*OrderCostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OrderCart_ServiceDesc is the grpc.ServiceDesc for OrderCart service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OrderCart_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ordercart.OrderCart",
	HandlerType: (*OrderCartServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOrderCost",
			Handler:    _OrderCart_GetOrderCost_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ordercart/ordercart.proto",
}

// CustomerNotificationClient is the client API for CustomerNotification service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CustomerNotificationClient interface {
}

type customerNotificationClient struct {
	cc grpc.ClientConnInterface
}

func NewCustomerNotificationClient(cc grpc.ClientConnInterface) CustomerNotificationClient {
	return &customerNotificationClient{cc}
}

// CustomerNotificationServer is the server API for CustomerNotification service.
// All implementations must embed UnimplementedCustomerNotificationServer
// for forward compatibility
type CustomerNotificationServer interface {
	mustEmbedUnimplementedCustomerNotificationServer()
}

// UnimplementedCustomerNotificationServer must be embedded to have forward compatible implementations.
type UnimplementedCustomerNotificationServer struct {
}

func (UnimplementedCustomerNotificationServer) mustEmbedUnimplementedCustomerNotificationServer() {}

// UnsafeCustomerNotificationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CustomerNotificationServer will
// result in compilation errors.
type UnsafeCustomerNotificationServer interface {
	mustEmbedUnimplementedCustomerNotificationServer()
}

func RegisterCustomerNotificationServer(s grpc.ServiceRegistrar, srv CustomerNotificationServer) {
	s.RegisterService(&CustomerNotification_ServiceDesc, srv)
}

// CustomerNotification_ServiceDesc is the grpc.ServiceDesc for CustomerNotification service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CustomerNotification_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ordercart.CustomerNotification",
	HandlerType: (*CustomerNotificationServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "ordercart/ordercart.proto",
}
