package middlewares

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func DbContext(client *mongo.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			req := c.Request()
			c.SetRequest(req.WithContext(
				context.WithValue(
					req.Context(),
					"CLIENT",
					client,
				),
			))

			if err := next(c); err != nil {
				return echo.NewHTTPError(500, err.Error())
			}

			return nil
		}
	}
}
