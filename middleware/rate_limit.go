/*******************************
 Name: ratelimit_middleware.go
 Date: 2022/12/12
 User: test
 Desc: 限流中间件
 Refer:
    -
*******************************/

package middleware

import (
	"strconv"
	"time"

	"github.com/secoba/echo"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

var (
	ipRateLimiter *limiter.Limiter
	store         limiter.Store
)

type (
	RateLimitConfig struct {
		Limit      int64
		StatusCode int
		StatusMsg  string
	}
)

var (
	DefaultRateLimitConfig = RateLimitConfig{
		Limit:      0,
		StatusCode: 403,
		StatusMsg:  "access too fast",
	}
)

func MiddlewareRateLimit() echo.MiddlewareFunc {
	c := DefaultRateLimitConfig
	return MiddlewareRateLimitWithConfig(c)
}

// MiddlewareRateLimitWithConfig https://github.com/ulule/limiter-examples
func MiddlewareRateLimitWithConfig(config RateLimitConfig) echo.MiddlewareFunc {
	// -------- no limit
	if config.Limit < 1 {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) (err error) {
				return next(c)
			}
		}
	}

	// -------- limit
	// 1. Configure
	rate := limiter.Rate{
		Period: time.Second,
		Limit:  config.Limit,
	}
	store = memory.NewStore()
	ipRateLimiter = limiter.New(store, rate)

	// 2. Return middleware handler
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			ip := c.RealIP()
			limiterCtx, err := ipRateLimiter.Get(c.Request().Context(), ip)
			if err != nil {
				//log(fmt.Sprintf("IPRateLimit - ipRateLimiter.Get - err: %v, %s on %s", err, ip, c.Request().URL), nil)
				return echo.NewHTTPError(config.StatusCode, config.StatusMsg)
			}

			h := c.Response().Header()
			h.Set("X-RateLimit-Limit", strconv.FormatInt(limiterCtx.Limit, 10))
			h.Set("X-RateLimit-Remaining", strconv.FormatInt(limiterCtx.Remaining, 10))
			h.Set("X-RateLimit-Reset", strconv.FormatInt(limiterCtx.Reset, 10))

			if limiterCtx.Reached {
				// logger.LogErr(fmt.Sprintf("Too Many Requests from %s on %s", ip, c.Request().URL))
				return echo.NewHTTPError(config.StatusCode, config.StatusMsg)
			}
			return next(c)
		}
	}
}
