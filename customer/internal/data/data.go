package data

import (
	"context"
	"customer/internal/conf"
	"fmt"
	"github.com/redis/go-redis/v9"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewCustomerRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	redis *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		redis: initRedis(c),
	}, cleanup, nil
}

func initRedis(c *conf.Data) *redis.Client {
	url := fmt.Sprintf("redis://%s/1?dial_timeout=1", c.Redis.GetAddr())
	options, _ := redis.ParseURL(url)
	client := redis.NewClient(options)
	status := client.Ping(context.Background())
	if status.Err() != nil {
		panic(status.Err())
	}
	return client
}
