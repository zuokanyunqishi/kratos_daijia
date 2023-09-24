package service

import (
	"context"

	pb "customer/api/customer"
)

type CustomerService struct {
	pb.UnimplementedCustomerServer
}

func NewCustomerService() *CustomerService {
	return &CustomerService{}
}

func (s *CustomerService) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.GetCustomerReply, error) {
	return &pb.GetCustomerReply{}, nil
}
