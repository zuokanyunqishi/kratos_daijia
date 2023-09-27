package middlewares

import (
	"context"
	"customer/internal/service"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/transport"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"strconv"
	"strings"
)

const (

	// bearerWord the bearer key word for authorization
	BearerWord string = "Bearer"

	// bearerFormat authorization token format
	bearerFormat string = "Bearer %s"

	// authorizationKey holds the key used to store the JWT Token in the request tokenHeader.
	AuthorizationKey string = "Authorization"

	// reason holds the error reason.
	Reason string = "UNAUTHORIZED"
)

func CustomerJwt(customerService *service.CustomerService) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {

			if header, ok := transport.FromServerContext(ctx); ok {

				token, ok := jwt.FromContext(ctx)
				if !ok {
					return nil, errors.Unauthorized(Reason, "claims not found")
				}
				claims := *(token.(*jwt2.MapClaims))
				customerId, err := strconv.ParseInt(claims["jti"].(string), 10, 64)
				if err != nil {
					return nil, errors.Unauthorized(Reason, "customerId invert err ")
				}
				auths := strings.SplitN(header.RequestHeader().Get(AuthorizationKey), " ", 2)
				jwtToken := auths[1]
				dbToken, err := customerService.GetTokenById(ctx, customerId)

				if err != nil || dbToken != jwtToken {
					return nil, errors.Unauthorized(Reason, "db token err ")
				}

				return handler(ctx, req)
			}

			return nil, nil
		}
	}

}
