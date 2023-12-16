package middlewares

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

func Cors() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				// Do something on entering
				tr.ReplyHeader().Set("Access-Control-Allow-Origin", "*")
				tr.ReplyHeader().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				tr.ReplyHeader().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
				tr.ReplyHeader().Set("Access-Control-Allow-Credentials", "true")
				defer func() {
					// Do something on exiting
				}()
			}
			// 跨域
			return handler(ctx, req)
		}
	}
}
