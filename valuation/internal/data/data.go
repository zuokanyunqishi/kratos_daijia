package data

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"valuation/internal/biz"
	"valuation/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewPrizeRuleData)

type Data struct {
	redis *redis.Client
	mysql *gorm.DB
}

func (d *Data) Redis() *redis.Client {
	return d.redis
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	data := &Data{
		redis: initRedis(c),
		mysql: initMysql(c),
	}
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		// 关掉 redis 连接
		data.redis.Close()
	}
	// 初始化表模型
	migrateTable(data.mysql)
	return data, cleanup, nil
}

func initRedis(c *conf.Data) *redis.Client {
	url := fmt.Sprintf("redis://%s/1?dial_timeout=%d", c.Redis.GetAddr(), 1)
	options, _ := redis.ParseURL(url)
	client := redis.NewClient(options)
	status := client.Ping(context.Background())
	if status.Err() != nil {
		panic(status.Err())
	}
	// Enable tracing
	if err := redisotel.InstrumentTracing(client); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	//if err := redisotel.InstrumentMetrics(client); err != nil {
	//	panic(err)
	//}
	return client
}

func initMysql(c *conf.Data) *gorm.DB {
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(err)
	}
	// gorm 接入 open tracing
	db.Use(tracing.NewPlugin(tracing.WithoutMetrics()))
	return db
}

func migrateTable(db *gorm.DB) {
	if err := db.AutoMigrate(&biz.PrizeRule{}); err != nil {
		log.Error(err)
	}

	rules := []biz.PrizeRule{
		{
			Model: gorm.Model{
				ID: 1,
			},
			PrizeRuleWork: biz.PrizeRuleWork{
				CityID:      1,
				StartFree:   300,
				DistanceFee: 35,
				DurationFee: 10,
				StartAt:     7,
				EndAt:       23,
			},
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			PrizeRuleWork: biz.PrizeRuleWork{
				CityID:      2,
				StartFree:   350,
				DistanceFee: 35,
				DurationFee: 10,
				StartAt:     23,
				EndAt:       24,
			},
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			PrizeRuleWork: biz.PrizeRuleWork{
				CityID:      3,
				StartFree:   400,
				DistanceFee: 35,
				DurationFee: 10,
				StartAt:     0,
				EndAt:       7,
			},
		},
	}
	db.Clauses(clause.OnConflict{DoNothing: true}).Create(rules)
}
