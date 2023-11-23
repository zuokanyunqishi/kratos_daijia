package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"valuation/internal/biz"
)

type PrizeRuleData struct {
	data *Data
	log  *log.Helper
}

func NewPrizeRuleData(data *Data, logger log.Logger) *PrizeRuleData {
	return &PrizeRuleData{data: data, log: log.NewHelper(logger)}
}

func (p *PrizeRuleData) GetRule(ctx context.Context, cityId uint, curr int) (*biz.PrizeRule, error) {
	//TODO implement me

	return nil, nil
}
