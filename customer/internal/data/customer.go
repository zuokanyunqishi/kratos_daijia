package data

import (
	"context"
	"customer/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type customerRepo struct {
	data *Data
	log  *log.Helper
}

func (r *customerRepo) SetPhoneCode(ctx context.Context, customer *biz.Customer) (*biz.Customer, error) {
	//TODO implement me
	r.data.redis.Set(ctx, "", "", 60)
	return nil, nil
}

// NewCustomerRepo
func NewCustomerRepo(data *Data, logger log.Logger) biz.CustomerRepo {
	return &customerRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
