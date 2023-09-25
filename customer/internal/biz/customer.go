package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type CustomerRepo interface {
	SetPhoneCode(ctx context.Context, customer *Customer) (*Customer, error)
}

// Customer Model
type Customer struct {
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

func (u *CustomerUsecase) SetCache(ctx context.Context, customer *Customer) {
	u.repo.SetPhoneCode(ctx, customer)
}
