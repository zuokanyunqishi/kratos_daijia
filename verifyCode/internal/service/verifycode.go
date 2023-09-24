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
		Code: randCode(req.Type, int(req.Length)),
	}, nil
}
func (s *VerifyCodeService) ListVerifyCode(ctx context.Context, req *pb.ListVerifyCodeRequest) (*pb.ListVerifyCodeReply, error) {
	return &pb.ListVerifyCodeReply{}, nil
}

type Chars struct {
	Value    string
	IdxBytes int
}

func randCode(t pb.TYPE, l int) string {
	chars := codeChars(t)
	return randCodeMixed(chars.Value, chars.IdxBytes, l)
}

func codeChars(t pb.TYPE) Chars {
	chars := Chars{}
	switch t {
	case pb.TYPE_Default:
	case pb.TYPE_DIGIT:
		chars.Value = "0123456789"
		//chars.IdxBytes = "1001"
		chars.IdxBytes = 0b1001
	case pb.TYPE_LETTER:
		chars.Value = "abcdefghijklmnopqrstuvwxyz"
		chars.IdxBytes = 0b11001
	case pb.TYPE_MIXED:
		chars.Value = "0123456789abcdefghijklmnopqrstuvwxyz"
		chars.IdxBytes = 0b100011
	}

	return chars
}
func randCodeSimple(chars string, l int) string {

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

func randCodeMixed(chars string, idxBits, l int) string {
	// 形成掩码
	idxMask := 1<<idxBits - 1
	// 63 位可以使用的最大组次数
	idxMax := 63 / idxBits

	// 利用string builder构建结果缓冲
	sb := strings.Builder{}
	sb.Grow(l)

	// 循环生成随机数
	// i 索引
	// cache 随机数缓存
	// remain 随机数还可以用几次
	for i, cache, remain := l-1, rand.Int63(), idxMax; i >= 0; {
		// 随机缓存不足，重新生成
		if remain == 0 {
			cache, remain = rand.Int63(), idxMax
		}
		// 利用掩码生成随机索引，有效索引为小于字符集合长度
		if idx := int(cache & int64(idxMask)); idx < len(chars) {
			sb.WriteByte(chars[idx])
			i--
		}
		// 利用下一组随机数位
		cache >>= idxBits
		remain--
	}

	return sb.String()
}
