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
	"strings"
)

type (
	VHostConfig struct {
		HostName   map[string]bool
		StatusCode int
		Msg        string
	}
)

var (
	DefaultVHostConfig = VHostConfig{
		HostName:   map[string]bool{},
		StatusCode: 403,
		Msg:        "cannot access",
	}
)

func MiddlewareVHost() echo.MiddlewareFunc {
	c := DefaultVHostConfig
	return MiddlewareVHostWithConfig(c)
}

func MiddlewareVHostWithConfig(config VHostConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			hostname := strings.SplitN(c.Request().Host, ":", 2)[0]
			if config.HostName != nil && len(config.HostName) > 0 {
				if _, ok := config.HostName[hostname]; ok {
					return next(c)
				}
				return echo.NewHTTPError(config.StatusCode, config.Msg)
			} else {
				return next(c)
			}
		}
	}
}
