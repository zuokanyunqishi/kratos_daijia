package data

import (
	"context"
	"driver/api/verifyCode"
	_ "driver/api/verifyCode"
	"driver/internal/biz"
	"driver/internal/conf"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"go.opentelemetry.io/otel"
	"time"
	_ "time"
)

type VerifyCode struct {
	*Data
	log *log.Helper
	rr  registry.Registrar
	cr  *conf.Registry
}

func (v *VerifyCode) GetVerifyCode(ctx context.Context, phone string, service string, lifeTime int64) (string, error) {

	ctx, span := otel.Tracer("data:driver").Start(ctx, "GetVerifyCode")
	defer span.End()
	v.log.WithContext(ctx).Infof("GetVerifyCode: %s, %s, %d", phone, service, lifeTime)

	code, err := v.makeVerifyCode(ctx, 6, verifyCode.TYPE_DIGIT)

	if err != nil {
		return "", err
	}

	ctx, span = otel.Tracer("data:driver:redis").Start(ctx, "GetVerifyCode")
	defer span.End()
	statusCmd := v.redis.Set(ctx, service+"verifyCode:"+phone, code, time.Second*time.Duration(lifeTime))
	if err = statusCmd.Err(); err != nil {
		return "", err
	}
	return code, nil

}

func (v *VerifyCode) ValidateVerifyCode(ctx context.Context, s string, s2 string) error {
	//TODO implement me
	panic("implement me")
}

func NewVerifyCode(data *Data, logger log.Logger, rr registry.Registrar, cr *conf.Registry) biz.VerifyCodeRepo {
	return &VerifyCode{Data: data,
		log: log.NewHelper(log.With(logger, "module", "data/verifyCode")),
		rr:  rr, cr: cr}
}

//func (v *VerifyCode) GetVerifyCode(ctx context.Context, phone string, service string, lifeTime int64) (string, error) {
//	ctx, span := otel.Tracer("data:driver").Start(ctx, "GetVerifyCode")
//	defer span.End()
//	v.log.WithContext(ctx).Infof("GetVerifyCode: %s, %s, %d", phone, service, lifeTime)
//
//	//code, err := v.makeVerifyCode(ctx, 6, verifyCode.TYPE_DIGIT)
//
//	//if err != nil {
//	//	return "", err
//	//}
//	code := "123456"
//	//statusCmd := v.redis.Set(ctx, service+"verifyCode:"+phone, code, time.Second*time.Duration(lifeTime))
//	//if err = statusCmd.Err(); err != nil {
//	//	return "", err
//	//}
//	return code, nil
//
//}
//
//func (v *VerifyCode) ValidateVerifyCode(ctx context.Context, phone string, service string) error {
//	//TODO implement me
//	panic("implement me")
//	return nil
//}

func (v *VerifyCode) makeVerifyCode(ctx context.Context, length uint32, t verifyCode.TYPE) (string, error) {

	//client, err := api.NewClient(api.DefaultConfig())
	//// new dis with consul client
	//dis := consul.New(client)
	endpoint := "discovery:///verifyCode"
	dis := v.rr.(*consul.Registry)
	conn, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint(endpoint), grpc.WithDiscovery(dis), grpc.WithMiddleware(tracing.Client()))

	if err != nil {
		return "", errors.New(555, "grpc init conn err", "")
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
