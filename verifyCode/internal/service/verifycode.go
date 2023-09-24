package service

import (
	"context"

	pb "verifyCode/api/verifyCode"
)

type VerifyCodeService struct {
	pb.UnimplementedVerifyCodeServer
}

func NewVerifyCodeService() *VerifyCodeService {
	return &VerifyCodeService{}
}

func (s *VerifyCodeService) CreateVerifyCode(ctx context.Context, req *pb.CreateVerifyCodeRequest) (*pb.CreateVerifyCodeReply, error) {
	return &pb.CreateVerifyCodeReply{}, nil
}
func (s *VerifyCodeService) UpdateVerifyCode(ctx context.Context, req *pb.UpdateVerifyCodeRequest) (*pb.UpdateVerifyCodeReply, error) {
	return &pb.UpdateVerifyCodeReply{}, nil
}
func (s *VerifyCodeService) DeleteVerifyCode(ctx context.Context, req *pb.DeleteVerifyCodeRequest) (*pb.DeleteVerifyCodeReply, error) {
	return &pb.DeleteVerifyCodeReply{}, nil
}
func (s *VerifyCodeService) GetVerifyCode(ctx context.Context, req *pb.GetVerifyCodeRequest) (*pb.GetVerifyCodeReply, error) {
	return &pb.GetVerifyCodeReply{}, nil
}
func (s *VerifyCodeService) ListVerifyCode(ctx context.Context, req *pb.ListVerifyCodeRequest) (*pb.ListVerifyCodeReply, error) {
	return &pb.ListVerifyCodeReply{}, nil
}
