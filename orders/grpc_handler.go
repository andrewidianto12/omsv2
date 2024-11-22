package main

import (
	"context"
	"log"

	pb "github.com/andrewidianto/common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer

	service OrderService
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrderService) {
	handler := &grpcHandler{
		service: service,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {

	log.Printf("new Order recevived %v", p)
	o := &pb.Order{
		ID: "42",
	}
	return o, nil

}
