package data

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"valuation/internal/biz"
	"valuation/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

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
	return db
}

func migrateTable(db *gorm.DB) {
	if err := db.AutoMigrate(&biz.PrizeRule{}); err != nil {
		log.Error(err)
	}
}
