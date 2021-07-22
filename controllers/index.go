package controllers

import (
	"github.com/labstack/echo/v4"
)

func AppRouter(e *echo.Echo) {
	e.GET("/search", Search)
}
