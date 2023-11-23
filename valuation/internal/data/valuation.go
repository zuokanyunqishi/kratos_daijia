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

func NewPrizeRuleData(data *Data, logger log.Logger) biz.PrizeRuleRepo {
	return &PrizeRuleData{data: data, log: log.NewHelper(logger)}
}

func (p *PrizeRuleData) GetRule(ctx context.Context, cityId uint, curr int) (*biz.PrizeRule, error) {
	pdata := &biz.PrizeRule{}

	result := p.data.mysql.Where("city_id = ? AND start_at >=  ? AND end_at < ?", cityId, curr, curr).WithContext(ctx).First(pdata)
	if result.Error != nil {
		return nil, result.Error
	}

	return pdata, nil
}
