package server

import (
	"context"
	"customer/api/customer"
	v1 "customer/api/helloworld/v1"
	"customer/internal/conf"
	"customer/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/transport/http"
	jwt2 "github.com/golang-jwt/jwt/v4"
)

func NewWhiteListMatcher() selector.MatchFunc {

	whiteList := make(map[string]struct{})
	whiteList["/shop.interface.v1.ShopInterface/Login"] = struct{}{}
	whiteList["/shop.interface.v1.ShopInterface/Register"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService,
	customerService *service.CustomerService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			selector.Server(
				jwt.Server(func(token *jwt2.Token) (interface{}, error) {
					return []byte("aaaa"), nil
				}, jwt.WithSigningMethod(jwt2.SigningMethodHS256), jwt.WithClaims(func() jwt2.Claims {
					return &jwt2.MapClaims{}
				})),
			).
				Match(NewWhiteListMatcher()).
				Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	// 注册customerService 服务
	customer.RegisterCustomerHTTPServer(srv, customerService)
	return srv
}
