package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/utils"
)

var excludedPaths = []string{
	"/v1",
	"/v1/",
}

func LoggerMiddleware(l *logrus.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) (err error) {
		t := time.Now()
		err = c.Next()

		if isExcludedPath(c.Path()) {
			return err
		}

		respBody := c.Response().Body()
		respStatus := c.Response().StatusCode()

		var bytes []byte
		bytes = append(bytes, respBody...)

		l.WithFields(logrus.Fields{
			"request":  getRequestLogFields(c),
			"response": getResponseLogFields(bytes, respStatus, t),
		}).Info("weblogger")

		return err
	}
}

func getRequestLogFields(c *fiber.Ctx) logrus.Fields {
	var bytes []byte

	reqBody := c.Request().Body()
	bytes = append(bytes, reqBody...)

	fields := logrus.Fields{
		"id":      c.Locals(utils.RequestIDKey),
		"method":  c.Method(),
		"path":    c.Path(),
		"url":     string(c.Request().RequestURI()),
		"headers": parseRequestHeaders(c.Request()),
		"body":    unmarshalBody(bytes),
	}

	return fields
}

func getResponseLogFields(body []byte, status int, t time.Time) logrus.Fields {
	return logrus.Fields{
		"status":   status,
		"duration": fmt.Sprint(time.Since(t).Round(time.Millisecond)),
		"body":     getResponseBody(body),
	}
}

func parseRequestHeaders(r *fasthttp.Request) map[string]interface{} {
	headers := make(map[string]interface{})

	r.Header.VisitAll(func(k, v []byte) {
		headers[string(k)] = string(v)
	})

	return headers
}

func unmarshalBody(b []byte) interface{} {
	var i interface{}
	_ = json.Unmarshal(b, &i)

	return i
}

func getResponseBody(b []byte) interface{} {
	if b == nil {
		return fiber.Map{"data": fiber.Map{}}
	}

	return unmarshalBody(b)
}

func isExcludedPath(path string) bool {
	for _, p := range excludedPaths {
		if path == p {
			return true
		}
	}

	return false
}
