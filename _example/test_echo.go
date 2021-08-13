package main

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/secoba/echo"
	"net/http"
)

func main() {
	// 创建 echo
	e := echo.New()
	g := e.Group("/api", "用户管理")
	g3 := e.Group("/test", "用户管理")
	e.GET("/home", "", func(context echo.Context) error {
		return context.HTML(http.StatusOK, "home")
	})
	g.GET("/test", "测试", func(context echo.Context) error {
		return context.HTML(http.StatusOK, "test")
	})
	g3.GET("/router", "", func(context echo.Context) error {
		return context.JSON(http.StatusOK, e.Routers())
	})
	g.GET("/routers", "", func(context echo.Context) error {
		return context.JSON(http.StatusOK, e.Routes())
	})
	g.GET("/group", "", func(context echo.Context) error {
		return context.JSON(http.StatusOK, e.RoutesGroup())
	})
	if err := e.Start(fmt.Sprintf(":%d", 18081)); err != nil {
		log.Error(err)
	}
}
