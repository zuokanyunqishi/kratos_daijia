package biz

import (
	"context"
	"customer/api/verifyCode"
	"github.com/go-kratos/kratos/v2/log"
)

type CustomerRepo interface {
	CachePhoneCode(ctx context.Context, customer *Customer, liftTime int64) error
	GetVerifyCode(ctx context.Context, length uint32, t verifyCode.TYPE) (string, error)
}

// Customer Model
type Customer struct {
	Telephone     string
	TelephoneCode string
}

// CustomerUsecase GreeterUsecase is a Customer usecase.
type CustomerUsecase struct {
	repo CustomerRepo
	log  *log.Helper
}

// NewCustomerUsecase NewGreeterUsecase new a Customer usecase.
func NewCustomerUsecase(repo CustomerRepo, logger log.Logger) *CustomerUsecase {
	return &CustomerUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (u *CustomerUsecase) SetVerifyCode(ctx context.Context, phone, code string, expireTime int64) error {
	return u.repo.CachePhoneCode(ctx, &Customer{
		Telephone:     phone,
		TelephoneCode: code,
	}, expireTime)
}

func (u *CustomerUsecase) GetVerifyCode(ctx context.Context, length uint32, t verifyCode.TYPE) (string, error) {
	return u.repo.GetVerifyCode(ctx, length, t)
}
