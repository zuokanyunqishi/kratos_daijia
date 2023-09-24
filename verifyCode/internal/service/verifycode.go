package service

import (
	"context"
	"math/rand"
	"strings"
	"time"

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
	return &pb.GetVerifyCodeReply{
		Code: randCode(codeChars(req.Type), int(req.Length)),
	}, nil
}
func (s *VerifyCodeService) ListVerifyCode(ctx context.Context, req *pb.ListVerifyCodeRequest) (*pb.ListVerifyCodeReply, error) {
	return &pb.ListVerifyCodeReply{}, nil
}

func codeChars(t pb.TYPE) string {
	var chars string
	switch t {
	case pb.TYPE_Default:
	case pb.TYPE_DIGIT:
		chars = "0123456789"
	case pb.TYPE_LETTER:
		chars = "abcdefghijklmnopqrstuvwxyz"
	case pb.TYPE_MIXED:
		chars = "0123456789abcdefghijklmnopqrstuvwxyz"
	}
	return chars
}
func randCode(chars string, l int) string {

	// 利用string builder构建结果缓冲
	sb := strings.Builder{}
	sb.Grow(l)
	charsLen := len(chars)

	rn := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		randIndex := rn.Intn(charsLen)
		sb.WriteByte(chars[randIndex])
	}

	return sb.String()
}
