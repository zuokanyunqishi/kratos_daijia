package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type CustomerRepo interface {
	CachePhoneCode(ctx context.Context, customer *Customer) error
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

func (u *CustomerUsecase) SetPhoneCode(ctx context.Context, phone, code string, expireTime int64) error {
	return u.repo.CachePhoneCode(ctx, &Customer{
		Telephone:     phone,
		TelephoneCode: code,
	})
}
