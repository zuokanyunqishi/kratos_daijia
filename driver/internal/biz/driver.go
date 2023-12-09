package biz

import (
	"context"
	"go.opentelemetry.io/otel"
)

type VerifyCodeRepo interface {
	GetVerifyCode(ctx context.Context, phone string, service string, lifeTime int64) (string, error)
	ValidateVerifyCode(ctx context.Context, phone string, service string) error
}

type DriverBiz struct {
	vc VerifyCodeRepo
}

func NewDriverBiz(vc VerifyCodeRepo) *DriverBiz {
	return &DriverBiz{vc: vc}
}

func (d *DriverBiz) GetVerifyCode(ctx context.Context, phone string, expireTime int64) (string, error) {
	ctx, span := otel.Tracer("biz:driver").Start(ctx, "GetVerifyCode")
	defer span.End()
	return d.vc.GetVerifyCode(ctx, phone, "driver:", expireTime)
}
