package biz

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"gorm.io/gorm"
	"valuation/api/mapService"
	"valuation/internal/conf"
)

type PrizeRule struct {
	gorm.Model
	PrizeRuleWork
}

type PrizeRuleWork struct {
	CityID      uint  `json:"city_id" gorm:"type:int;not null"`
	StartFree   int64 `json:"start_free" gorm:"type:int;not null"`
	DistanceFee int64 `json:"distance_free" gorm:"type:int;not null"`
	DurationFee int64 `json:"duration_fee" gorm:"type:int;not null"`
	StartAt     int   `json:"start_at" gorm:"type:int;not null"`
	EndAt       int   `json:"end_at" gorm:"type:int;not null"`
}

type PrizeRuleRepo interface {
	GetRule(ctx context.Context, cityId uint, curr int) (*PrizeRule, error)
}

type ValuationBiz struct {
	log *log.Helper
	pri PrizeRuleRepo
	rr  registry.Registrar
	cr  *conf.Registry
}

func NewValuationBiz(logger log.Logger, pri PrizeRuleRepo, rr registry.Registrar, cr *conf.Registry) *ValuationBiz {
	return &ValuationBiz{log: log.NewHelper(logger), pri: pri, rr: rr, cr: cr}
}

func (b *ValuationBiz) GetRuleInfo(ctx context.Context, cityId uint, curr int) (*PrizeRule, error) {

	return b.pri.GetRule(ctx, cityId, curr)
}

func (b *ValuationBiz) GetDrivingInfo(ctx context.Context, origin, destination string) (string, error) {

	endpoint := "discovery:///verifyCode"
	dis := b.rr.(*consul.Registry)
	conn, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint(endpoint), grpc.WithDiscovery(dis))

	if err != nil {
		return "", errors.New("grpc init conn err")
	}
	defer conn.Close()
	client := mapService.NewMapServiceClient(conn)
	info, err := client.GetDriverInfo(ctx, &mapService.GetDriverInfoReq{
		Origin:      origin,
		Destination: destination,
	})

	fmt.Println(info)
	return "", err
}
