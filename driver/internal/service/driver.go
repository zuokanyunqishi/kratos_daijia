package service

import (
	"context"
	pb "driver/api/driver"
	"driver/internal/biz"
	"time"
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
		return &pb.GetVerifyCoderRes{
			Message: err.Error(),
			Code:    500,
		}, err
	}
	return &pb.GetVerifyCoderRes{VerifyCode: code,
		Message:        "验证码发送成功",
		Code:           200,
		VerifyCodeLife: time.Now().Add(600 * time.Second).Unix(),
		VerifyCodeTime: time.Now().Unix(),
	}, nil
}
