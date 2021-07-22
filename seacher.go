package main

import (
	"net/http"

	"prsSearcher/configs"
	"prsSearcher/controllers"
	"prsSearcher/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	// DB 접속
	client, ctx, cancel := configs.ConnectDB()

	// 함수 종료 뒤 연결을 끊어지도록 설정
	defer client.Disconnect(ctx)
	defer cancel()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://192.168.160.15:80", "http://localhost:8080"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	// 각 request마다 고유의 ID를 부여
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middlewares.DbContext(client))
	e.Use(middlewares.LogrusLogger())

	controllers.AppRouter(e)

	// e.Logger.Fatal(e.Start(":1324"))
	e.Logger.Fatal(e.Start("192.168.160.15:8088"))
}
