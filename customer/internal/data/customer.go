package data

import (
	"context"
	"customer/api/verifyCode"
	"customer/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"time"
)

type customerData struct {
	data *Data
	log  *log.Helper
}

func (r *customerData) GetVerifyCode(ctx context.Context, length uint32, t verifyCode.TYPE) (string, error) {
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

func (r *customerData) CachePhoneCode(ctx context.Context, phone, verifyCode string, lifeTime int64) error {
	statusCmd := r.data.redis.Set(ctx, "CachePhoneCode:"+phone, verifyCode, time.Second*time.Duration(lifeTime))
	return statusCmd.Err()
}

func NewCustomerData(data *Data, logger log.Logger) biz.CustomerRepo {
	return &customerData{
		data: data,
		log:  log.NewHelper(logger),
	}
}
