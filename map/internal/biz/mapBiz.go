package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type MapInfoRepo interface {
	GetDrivingInfo(ctx context.Context, origin, destination string)
}

type MapBizUsecase struct {
	repo *MapInfoRepo
	log  *log.Helper
}
