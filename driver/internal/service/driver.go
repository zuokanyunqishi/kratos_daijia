package service

import (
	"context"
	pb "driver/api/driver"
	"driver/internal/biz"
)

type DriverService struct {
	pb.UnimplementedDriverServer
	driverBiz *biz.DriverBiz
}

func NewDriverService(driverBiz *biz.DriverBiz) *DriverService {
	return &DriverService{driverBiz: driverBiz}
}

func (s *DriverService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeReq) (*pb.GetVerifyCoderRes, error) {
	code, err := s.driverBiz.GetVerifyCode(ctx, req.Telephone, 600)
	if err != nil {
		return nil, err
	}
	return &pb.GetVerifyCoderRes{VerifyCode: code}, nil
}
