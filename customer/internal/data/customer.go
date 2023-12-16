package data

import (
	"context"
	"customer/api/verifyCode"
	"customer/internal/biz"
	"customer/internal/conf"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"time"
)

type customerData struct {
	data *Data
	log  *log.Helper
	rr   registry.Registrar
	cr   *conf.Registry
}

func (d *customerData) DeleteToken(ctx context.Context, id int64) error {
	var customer biz.Customer
	customer.ID = uint(id)
	customer.Token = ""
	result := d.data.mysql.WithContext(ctx).
		Model(&customer).
		Select("token").
		Updates(customer)
	return result.Error
}

func (d *customerData) GetTokenById(ctx context.Context, id int64) (string, error) {

	var customer biz.Customer
	result := d.data.mysql.WithContext(ctx).Select("token").First(&customer, "id = ?", id)

	if result.Error != nil || result.RowsAffected <= 0 {
		return "", errors.New("not found token")
	}
	return customer.Token, nil
}

func (d *customerData) UpdateCustomerToken(ctx context.Context, c *biz.Customer) (*biz.Customer, error) {
	var customer biz.Customer
	customer.ID = c.ID
	customer.Token = c.Token
	customer.TokenCreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	//customer.UpdatedAt = time.Now()
	result := d.data.mysql.
		WithContext(ctx).
		Model(&customer).
		Select("token", "token_created_at").
		Updates(customer)
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

func NewCustomerData(data *Data, logger log.Logger, rr registry.Registrar, cr *conf.Registry) biz.CustomerRepo {
	return &customerData{
		data: data,
		log:  log.NewHelper(logger),
		rr:   rr,
		cr:   cr,
	}
}

func (d *customerData) MakeVerifyCode(ctx context.Context, length uint32, t verifyCode.TYPE) (string, error) {

	//client, err := api.NewClient(api.DefaultConfig())
	//// new dis with consul client
	//dis := consul.New(client)
	endpoint := "discovery:///verifyCode"
	dis := d.rr.(*consul.Registry)
	conn, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint(endpoint),
		grpc.WithDiscovery(dis),
		grpc.WithMiddleware(tracing.Client()))

	if err != nil {
		return "", errors.New("grpc init conn err")
	}
	defer conn.Close()
	// 构建客户端
	verifyCodeClient := verifyCode.NewVerifyCodeClient(conn)
	code, err := verifyCodeClient.GetVerifyCode(ctx, &verifyCode.GetVerifyCodeRequest{
		Length: length,
		Type:   t,
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
	return d.data.redis.Get(ctx, "CachePhoneCode:"+telephone).Val()
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
