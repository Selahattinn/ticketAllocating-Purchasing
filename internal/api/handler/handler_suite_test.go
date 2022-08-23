//go:build integration
// +build integration

package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus/hooks/test"

	di "github.com/Selahattinn/ticketAllocating-Purchasing"
	"github.com/Selahattinn/ticketAllocating-Purchasing/configs"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/constants"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/healthcheck"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/mysql"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/response"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/utils"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/validation"
)

const (
	reqContentTypeHeaderKey = "Content-Type"
	reqJSONHeaderValue      = "application/json"

	emptyPayload = `{ "data": {} }`
)

var (
	ctx           context.Context
	req           *http.Request
	server        *fiber.App
	mockCtrl      *gomock.Controller
	mysqlInstance mysql.IMysqlInstance
)

func TestAPIIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "API Integration Suite")
}

var _ = BeforeSuite(func() {
	var err error

	// Setup instance
	mysqlInstance, err = mysql.InitMysql(mysql.Config{
		URL: fmt.Sprintf("%s", constants.MysqlTestURL),
	})

	Expect(err).NotTo(HaveOccurred())

	// Setup context & mock dependencies
	loadMockDependencies()

	// Setup application
	initApplication()

	go func() {
		if serveErr := server.Listen(fmt.Sprintf(":%s", constants.ServerTestPort)); serveErr != nil {
			Expect(serveErr).NotTo(HaveOccurred())
		}
	}()
})

var _ = AfterSuite(func() {
	defer func() {
		_ = mysqlInstance.Database().Close()
		_ = server.Shutdown()
	}()
})

func loadMockDependencies() {
	mockCtrl = gomock.NewController(GinkgoT())

	configs.TicketApp = &configs.TicketScheme{
		Web: configs.WebConfig{
			Env: constants.AppTestEnv,
		},
	}
}

func initApplication() {
	logger, _ := test.NewNullLogger()
	server = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	healthcheck.InitHealthCheck()
	validator := validation.InitValidator()

	ctx = context.WithValue(context.Background(), utils.ValidatorKey, validator) //nolint:staticcheck

	server.Use(func(c *fiber.Ctx) error {
		c.Locals(utils.ContextKey, ctx)
		return c.Next()
	})

	route := di.InitRoute(
		logger,
		mysqlInstance,
	)

	route.SetupRoutes(&api.RouteContext{
		App: server,
	})
}

func prepareRequest(method, url string, data []byte) {
	var body io.Reader
	if data != nil {
		body = bytes.NewBuffer(data)
	}

	req = httptest.NewRequest(method, url, body)
	req.Header.Add(reqContentTypeHeaderKey, reqJSONHeaderValue)
}

func sendTestRequest(req *http.Request) (*http.Response, []byte) {
	resp, _ := server.Test(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return resp, body
}

func assertErrorResponseByHTTPResponse(resp *http.Response, data []byte, code string, statusCode int) {
	var httpErrResp response.HTTPErrorResponse
	actual := json.Unmarshal(data, &httpErrResp)

	Expect(actual).NotTo(HaveOccurred())
	Expect(httpErrResp.Error.Code).To(Equal(code))
	Expect(resp.StatusCode).To(Equal(statusCode))
}

func assertBodyParserFailed(resp *http.Response, data []byte) {
	assertErrorResponseByHTTPResponse(resp, data, "body_parser_failed", fiber.StatusUnprocessableEntity)
}

func assertValidationFailed(resp *http.Response, data []byte) {
	assertErrorResponseByHTTPResponse(resp, data, "validation_failed", fiber.StatusUnprocessableEntity)
}

func assertErrorResponseByCodeWithMessage(data []byte, code string, message string) {
	Expect(data).To(MatchJSON(fmt.Sprintf(`{"error":{"code":"%s","message":"%s"}}`, code, message)))
}

func prepareBody(i interface{}) []byte {
	payload, _ := json.Marshal(i)
	return payload
}
