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
	AllowRemoteConfig struct {
		Address map[string]bool
	}
)

var (
	DefaultAllowRemoteConfig = AllowRemoteConfig{
		Address: map[string]bool{},
	}
)

func MiddlewareAllowRemote() echo.MiddlewareFunc {
	c := DefaultAllowRemoteConfig
	return MiddlewareAllowRemoteWithConfig(c)
}

func MiddlewareAllowRemoteWithConfig(config AllowRemoteConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			remote := c.RealIP()
			if config.Address != nil && len(config.Address) > 0 {
				if _, ok := config.Address[remote]; ok {
					return next(c)
				}
				return echo.NewHTTPError(http.StatusForbidden, "cannot access")
			} else {
				return next(c)
			}
		}
	}
}
