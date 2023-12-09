package data

import (
	"context"
	"driver/internal/biz"
	"driver/internal/conf"
	"fmt"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewVerifyCode)

// Data .
type Data struct {
	// TODO wrapped database client
	redis redis.UniversalClient
	mysql *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {

	data := &Data{
		redis: initRedis(c),
		mysql: initMysql(c),
	}
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	migrateTable(data.mysql)
	return data, cleanup, nil
}

func migrateTable(db *gorm.DB) {
	if err := db.AutoMigrate(
		&biz.Driver{},
	); err != nil {
		panic(err)
	}
}

func initRedis(c *conf.Data) *redis.Client {
	url := fmt.Sprintf("redis://%s/1?dial_timeout=%d", c.Redis.GetAddr(), 1)
	options, _ := redis.ParseURL(url)
	client := redis.NewClient(options)
	status := client.Ping(context.Background())
	if status.Err() != nil {
		panic(status.Err())
	}
	// Enable tracing instrumentation.
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
