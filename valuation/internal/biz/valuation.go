package biz

import (
	"gorm.io/gorm"
	"time"
)

type PrizeRule struct {
	gorm.Model
	PrizeRuleWork
}

type PrizeRuleWork struct {
	CityID      uint      `json:"city_id" gorm:"type:int;not null"`
	StartFree   int64     `json:"start_free" gorm:"type:int;not null"`
	DistanceFee int64     `json:"distance_free" gorm:"type:int;not null"`
	DurationFee int64     `json:"duration_fee" gorm:"type:int;not null"`
	StartAt     time.Time `json:"start_at" gorm:"type:time;not null"`
	EndAt       time.Time `json:"end_at" gorm:"type:time;not null"`
}
