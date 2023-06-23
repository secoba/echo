package middleware

import (
	"net/http"

	"github.com/secoba/echo"
)

// RedirectConfig defines the config for Redirect middleware.
type RedirectConfig struct {
	// Skipper defines a function to skip middleware.
	Skipper

	// Status code to be used when redirecting the request.
	// Optional. Default value http.StatusMovedPermanently.
	Code int `yaml:"code"`
}

// redirectLogic represents a function that given a scheme, host and uri
// can both: 1) determine if redirect is needed (will set ok accordingly) and
// 2) return the appropriate redirect url.
type redirectLogic func(scheme, host, uri string) (ok bool, url string)

const www = "www."

// DefaultRedirectConfig is the default Redirect middleware config.
var DefaultRedirectConfig = RedirectConfig{
	Skipper: DefaultSkipper,
	Code:    http.StatusMovedPermanently,
}

// HTTPSRedirect redirects http requests to https.
// For example, http://labstack.com will be redirect to https://labstack.com.
//
// Usage `Echo#Pre(HTTPSRedirect())`
func MiddlewareHTTPSRedirect() echo.MiddlewareFunc {
	return MiddlewareHTTPSRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `HTTPSRedirect()`.
func MiddlewareHTTPSRedirectWithConfig(config RedirectConfig) echo.MiddlewareFunc {
	return redirect(config, func(scheme, host, uri string) (ok bool, url string) {
		if ok = scheme != "https"; ok {
			url = "https://" + host + uri
		}
		return
	})
}

// HTTPSWWWRedirect redirects http requests to https www.
// For example, http://labstack.com will be redirect to https://www.labstack.com.
//
// Usage `Echo#Pre(HTTPSWWWRedirect())`
func MiddlewareHTTPSWWWRedirect() echo.MiddlewareFunc {
	return MiddlewareHTTPSWWWRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSWWWRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `HTTPSWWWRedirect()`.
func MiddlewareHTTPSWWWRedirectWithConfig(config RedirectConfig) echo.MiddlewareFunc {
	return redirect(config, func(scheme, host, uri string) (ok bool, url string) {
		if ok = scheme != "https" && host[:4] != www; ok {
			url = "https://www." + host + uri
		}
		return
	})
}

// HTTPSNonWWWRedirect redirects http requests to https non www.
// For example, http://www.labstack.com will be redirect to https://labstack.com.
//
// Usage `Echo#Pre(HTTPSNonWWWRedirect())`
func MiddlewareHTTPSNonWWWRedirect() echo.MiddlewareFunc {
	return MiddlewareHTTPSNonWWWRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSNonWWWRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `HTTPSNonWWWRedirect()`.
func MiddlewareHTTPSNonWWWRedirectWithConfig(config RedirectConfig) echo.MiddlewareFunc {
	return redirect(config, func(scheme, host, uri string) (ok bool, url string) {
		if ok = scheme != "https"; ok {
			if host[:4] == www {
				host = host[4:]
			}
			url = "https://" + host + uri
		}
		return
	})
}

// WWWRedirect redirects non www requests to www.
// For example, http://labstack.com will be redirect to http://www.labstack.com.
//
// Usage `Echo#Pre(WWWRedirect())`
func MiddlewareWWWRedirect() echo.MiddlewareFunc {
	return MiddlewareWWWRedirectWithConfig(DefaultRedirectConfig)
}

// WWWRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `WWWRedirect()`.
func MiddlewareWWWRedirectWithConfig(config RedirectConfig) echo.MiddlewareFunc {
	return redirect(config, func(scheme, host, uri string) (ok bool, url string) {
		if ok = host[:4] != www; ok {
			url = scheme + "://www." + host + uri
		}
		return
	})
}

// NonWWWRedirect redirects www requests to non www.
// For example, http://www.labstack.com will be redirect to http://labstack.com.
//
// Usage `Echo#Pre(NonWWWRedirect())`
func MiddlewareNonWWWRedirect() echo.MiddlewareFunc {
	return MiddlewareNonWWWRedirectWithConfig(DefaultRedirectConfig)
}

// NonWWWRedirectWithConfig returns an HTTPSRedirect middleware with config.
// See `NonWWWRedirect()`.
func MiddlewareNonWWWRedirectWithConfig(config RedirectConfig) echo.MiddlewareFunc {
	return redirect(config, func(scheme, host, uri string) (ok bool, url string) {
		if ok = host[:4] == www; ok {
			url = scheme + "://" + host[4:] + uri
		}
		return
	})
}

func redirect(config RedirectConfig, cb redirectLogic) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultTrailingSlashConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			req, scheme := c.Request(), c.Scheme()
			host := req.Host
			if ok, url := cb(scheme, host, req.RequestURI); ok {
				return c.Redirect(config.Code, url)
			}

			return next(c)
		}
	}
}
