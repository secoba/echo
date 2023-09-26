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
	"net/http"

	"github.com/secoba/echo"
)

type (
	VHostConfig struct {
		HostName map[string]bool
	}
)

var (
	DefaultVHostConfig = VHostConfig{
		HostName: map[string]bool{},
	}
)

func MiddlewareVHost() echo.MiddlewareFunc {
	c := DefaultVHostConfig
	return MiddlewareVHostWithConfig(c)
}

func MiddlewareVHostWithConfig(config VHostConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			hostname := c.Request().Host
			if config.HostName != nil && len(config.HostName) > 0 {
				if _, ok := config.HostName[hostname]; ok {
					return next(c)
				}
				return echo.NewHTTPError(http.StatusForbidden, "cannot access")
			} else {
				return next(c)
			}
		}
	}
}
