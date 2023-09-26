/*******************************
 Name: acl_middleware.go
 Date: 2022/12/12
 User: test
 Desc: 访问控制中间件
 Refer:
    -
*******************************/

package middleware

import (
	"github.com/secoba/echo"
)

type (
	DynamicHeadersConfig struct {
		Headers map[string]func(ctx echo.Context) string
	}
)

var (
	DefaultDynamicHeadersConfig = DynamicHeadersConfig{
		Headers: nil,
	}
)

func MiddlewareDynamicHeaders() echo.MiddlewareFunc {
	c := DefaultDynamicHeadersConfig
	return MiddlewareDynamicHeadersWithConfig(c)
}

func MiddlewareDynamicHeadersWithConfig(config DynamicHeadersConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Headers != nil {
				for key, call := range config.Headers {
					c.Request().Header.Set(key, call(c))
				}
			}
			return next(c)
		}
	}
}
