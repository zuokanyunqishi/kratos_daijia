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

func (r *customerRepo) CachePhoneCode(ctx context.Context, customer *biz.Customer) error {

	statusCmd := r.data.redis.Set(ctx, "CachePhoneCode:"+customer.Telephone, customer.TelephoneCode, 60)
	return statusCmd.Err()
}

func NewCustomerRepo(data *Data, logger log.Logger) biz.CustomerRepo {
	return &customerRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
