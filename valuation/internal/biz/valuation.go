package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
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
}

func (b *ValuationBiz) GetRuleInfo(ctx context.Context, cityId uint, curr int) (*PrizeRule, error) {

	return b.pri.GetRule(ctx, cityId, curr)
}

func NewValuationBiz(pri PrizeRuleRepo, logger log.Logger) *ValuationBiz {
	return &ValuationBiz{log: log.NewHelper(logger), pri: pri}
}
