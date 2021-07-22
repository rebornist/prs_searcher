package middlewares

import (
	"context"
	"fmt"
	"os"
	"prsSearcher/models"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LogrusLogger() echo.MiddlewareFunc {
	/* ... logger 초기화 */
	logger := logrus.New()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logEntry := logrus.NewEntry(logger)
			// var logResponse config.Logger
			data := make(map[string]interface{})

			// var httpBody *http.body

			// request_id를 가져와 logEntry에 셋팅
			id := c.Request().Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = c.Response().Header().Get(echo.HeaderXRequestID)
			}

			var getBodyData []string
			fps, _ := c.FormParams()
			for k, v := range fps {
				value := fmt.Sprintf("%s: %s", k, strings.Join(v, "&"))
				getBodyData = append(getBodyData, value)
			}

			// logrus에 저장
			data["request_id"] = id
			data["body"] = getBodyData
			data["connect_ip"] = c.RealIP()
			data["request_url"] = c.Request().URL.RequestURI()
			data["user_agent"] = c.Request().UserAgent()

			logEntry = logEntry.WithFields(data)
			// logEntry를 Context에 저장
			req := c.Request()
			c.SetRequest(req.WithContext(
				context.WithValue(
					req.Context(),
					"LOG",
					logEntry,
				),
			))

			return next(c)
		}
	}
}

func CreateLogger(logger *logrus.Entry, status int, msg string) {

	logn := logrus.New()

	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// You could set this to any `io.Writer` such as a file
	_, week := time.Now().ISOWeek()
	file, err := os.OpenFile(fmt.Sprintf("log/log_%d.log", week), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logn.Out = file
	} else {
		logn.Error("Failed to log to file, using default stderr")
	}

	logConf := models.Logger{
		Body:       fmt.Sprintf("%s", logger.Data["body"]),
		ConnectIp:  fmt.Sprintf("%s", logger.Data["connect_ip"]),
		RequestId:  fmt.Sprintf("%s", logger.Data["request_id"]),
		RequestUrl: fmt.Sprintf("%s", logger.Data["request_url"]),
		Status:     status,
		Backoff:    time.Second.Milliseconds(),
		UserAgent:  fmt.Sprintf("%s", logger.Data["user_agent"]),
		CreatedAt:  time.Now(),
	}

	logs := logn.WithFields(logrus.Fields{
		"backoff":     logConf.Backoff,
		"body":        logConf.Body,
		"created":     logConf.CreatedAt,
		"IP":          logConf.ConnectIp,
		"request-id":  logConf.RequestId,
		"request-url": logConf.RequestUrl,
		"status":      logConf.Status,
		"user-agent":  logConf.UserAgent,
	})

	logs.Info(msg)
}
