package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		entry := logger.WithFields(logrus.Fields{
			"status_code":  c.Writer.Status(),
			"latency_time": latency,
			"client_ip":    c.ClientIP(),
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info("Request completed")
		}
	}
}
