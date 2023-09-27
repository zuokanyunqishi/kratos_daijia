package service

import (
	"context"
	"customer/api/verifyCode"
	"customer/internal/biz"
	"regexp"
	"time"

	pb "customer/api/customer"
)

type CustomerService struct {
	pb.UnimplementedCustomerServer
	cus *biz.CustomerUsecase
}

const TokenLifeTime = 60 * 60 * 24 * 30 * 2

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
	code, err := s.cus.MakeVerifyCode(ctx, 6, verifyCode.TYPE_DIGIT)
	if err != nil {
		return &pb.GetCustomerReply{
			Code:    1,
			Message: "验证码获取错误",
		}, nil
	}
	err = s.cus.CachePhoneCode(ctx, req.Telephone, code, 60)
	if err != nil {
		return &pb.GetCustomerReply{
			Code:    1,
			Message: "验证码缓存错误",
		}, nil
	}
	return &pb.GetCustomerReply{
		Code:           0,
		VerifyCode:     code,
		VerifyCodeTime: time.Now().Unix(),
		VerifyCodeLife: 60,
	}, nil
}

// Login login
func (s *CustomerService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginRes, error) {

	var customer biz.Customer
	var err error
	// 校验电话,和验证码
	code := s.cus.GetVerifyCode(ctx, req.Telephone)
	if code != req.VerifyCode {
		return &pb.LoginRes{
			Code:    1,
			Message: "验证码错误",
		}, nil
	}

	// 判断电话号码是否注册,获取客户信息
	if customer, err = s.cus.GetRepo().GetCustomerByTelephone(ctx, req.Telephone); err != nil {
		if customer, err = s.cus.GetRepo().QuickCreateCustomerByPhone(ctx, req.Telephone); err != nil {
			return &pb.LoginRes{
				Code:    1,
				Message: "创建用户错误 ",
			}, nil
		}
	}
	token, err := s.cus.GenerateTokenAndSave(ctx, &customer, TokenLifeTime*time.Second)
	if err != nil {
		return &pb.LoginRes{
			Code:    1,
			Message: "生成token失败",
		}, nil
	}

	// 生成token 返回数据
	return &pb.LoginRes{
		Code:          0,
		Message:       "登录成功",
		Token:         token,
		TokenCreateAt: time.Now().Unix(),
		TokenLifeTime: TokenLifeTime,
	}, nil
}

func (s *CustomerService) Logout(ctx context.Context, req *pb.LogoutReq) (*pb.LogoutRes, error) {

	return &pb.LogoutRes{}, nil
}

func (s *CustomerService) GetTokenById(ctx context.Context, id int64) (string, error) {
	return s.cus.GetRepo().GetTokenById(ctx, id)
}
