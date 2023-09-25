package service

import (
	"context"
	"customer/api/verifyCode"
	"customer/internal/biz"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"regexp"
	"time"

	pb "customer/api/customer"
)

type CustomerService struct {
	pb.UnimplementedCustomerServer
	cus *biz.CustomerUsecase
}

func NewCustomerService(cus *biz.CustomerUsecase) *CustomerService {
	return &CustomerService{
		cus: cus,
	}
}

func (s *CustomerService) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.GetCustomerReply, error) {
	//telephone := req.GetTelephone()
	//

	pattern := "^(13[0-9]|14[579]|15[0-3,5-9]|16[6]|17[0135678]|18[0-9]|19[89])\\d{8}$"
	compile := regexp.MustCompile(pattern)
	if !compile.MatchString(req.Telephone) {
		return &pb.GetCustomerReply{
			Code:    1,
			Message: "电话号码格式错误",
		}, nil
	}
	// 获取验证码
	//
	conn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint("localhost:9000"))
	defer conn.Close()

	// 构建客户端
	client := verifyCode.NewVerifyCodeClient(conn)
	code, err := client.GetVerifyCode(ctx, &verifyCode.GetVerifyCodeRequest{
		Length: 6,
		Type:   verifyCode.TYPE_DIGIT,
	})

	if err != nil {
		return &pb.GetCustomerReply{
			Code:    1,
			Message: "验证码获取错误",
		}, nil
	}
	err = s.cus.SetVerifyCode(ctx, req.Telephone, code.Code, 60)
	if err != nil {
		return &pb.GetCustomerReply{
			Code:    1,
			Message: "验证码缓存错误",
		}, nil
	}

	return &pb.GetCustomerReply{
		Code:           0,
		VerifyCode:     code.Code,
		VerifyCodeTime: time.Now().Unix(),
		VerifyCodeLife: 60,
	}, nil
}
