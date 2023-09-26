package data

import (
	"context"
	"customer/api/verifyCode"
	"customer/internal/biz"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"time"
)

type customerData struct {
	data *Data
	log  *log.Helper
}

func (d *customerData) UpdateCustomer(ctx context.Context, c *biz.Customer) (*biz.Customer, error) {
	result := d.data.mysql.WithContext(ctx).Save(c)
	return c, result.Error
}

func (d *customerData) QuickCreateCustomerByPhone(ctx context.Context, telephone string) (biz.Customer, error) {
	var customer biz.Customer
	customer.CustomerWork.Telephone = telephone
	result := d.data.mysql.WithContext(ctx).Create(&customer)

	if result.RowsAffected <= 0 {
		return customer, fmt.Errorf("%q,%w", "data:customer:QuickCreateCustomerByPhone createCustomer err", result.Error)
	}
	return customer, nil
}

func NewCustomerData(data *Data, logger log.Logger) biz.CustomerRepo {
	return &customerData{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (d *customerData) MakeVerifyCode(ctx context.Context, length uint32, t verifyCode.TYPE) (string, error) {
	conn, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint("localhost:9000"))
	defer conn.Close()
	// 构建客户端
	client := verifyCode.NewVerifyCodeClient(conn)
	code, err := client.GetVerifyCode(ctx, &verifyCode.GetVerifyCodeRequest{
		Length: length,
		Type:   verifyCode.TYPE_DIGIT,
	})

	if err != nil {
		return "", err
	}
	return code.Code, nil
}

func (d *customerData) CachePhoneCode(ctx context.Context, phone, verifyCode string, lifeTime int64) error {
	statusCmd := d.data.redis.Set(ctx, "CachePhoneCode:"+phone, verifyCode, time.Second*time.Duration(lifeTime))
	return statusCmd.Err()
}

func (d *customerData) GetVerifyCode(ctx context.Context, telephone string) string {
	return d.data.redis.Get(ctx, "CachePhoneCode:"+telephone).String()
}

func (d *customerData) GetCustomerByTelephone(ctx context.Context, telephone string) (biz.Customer, error) {
	var customer biz.Customer
	result := d.data.mysql.WithContext(ctx).First(&customer, "telephone = ?", telephone)
	// 没有找到
	if result.RowsAffected <= 0 {
		return customer, errors.New("data:customer not found")
	}

	return customer, nil
}
