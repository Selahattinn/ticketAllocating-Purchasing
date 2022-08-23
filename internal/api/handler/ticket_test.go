//go:build integration
// +build integration

package handler_test

import (
	"fmt"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/dto/request"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strconv"
)

const (
	createTicketURL   = "/v1/ticket_options"
	purcahseTicketURL = "/v1/ticket_options/%s/purchase"
	getTicketURL      = "/v1/ticket/%s"
)

type ctStruct struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	Allocation  int    `json:"allocation"`
}

var _ = Describe("AccountIntegration", func() {
	BeforeEach(func() {
		_, _ = mysqlInstance.Database().Exec("DELETE FROM ticket_options")
		_, _ = mysqlInstance.Database().Exec("DELETE FROM ticket")

	})

	Describe("POST /ticket_options", func() {
		var (
			ctReq request.CreateTicketRequest
		)

		BeforeEach(func() {
			ctReq = request.CreateTicketRequest{
				Name:        gofakeit.Name(),
				Description: gofakeit.SentenceSimple(),
				Allocation:  gofakeit.IntRange(1, 100),
			}
		})

		It("should return create error when request is not parsed", func() {
			prepareRequest(fiber.MethodPost, createTicketURL, []byte(`{
			  "name": "test",
			  "desc": "sample test",
			  "allocation": "-1"
			}`))

			actual, body := sendTestRequest(req)

			assertBodyParserFailed(actual, body)
		})

		It("should return create error when request is not validated", func() {
			prepareRequest(fiber.MethodPost, createTicketURL, []byte(`{
			  "name": "test",
			  "desc": "sample test",
			  "allocation": -1
			}`))

			actual, body := sendTestRequest(req)

			assertValidationFailed(actual, body)
		})

		It("should return create ticket response", func() {
			prepareRequest(fiber.MethodPost, createTicketURL, prepareBody(ctReq))

			actual, body := sendTestRequest(req)
			rows, err := mysqlInstance.Database().Query("select * from ticket")

			Expect(err).To(BeNil())

			var ctRes ctStruct
			for rows.Next() {
				err := rows.Scan(&ctRes.ID, &ctRes.Name, &ctRes.Description, &ctRes.Allocation)
				Expect(err).ToNot(HaveOccurred())
			}

			expected := fmt.Sprintf(`{
			  "data": {
				"id": %d,
				"name": "%s",
				"desc": "%s",
				"allocation": %d
			  }
			}`, ctRes.ID, ctRes.Name, ctRes.Description, ctRes.Allocation)

			Expect(actual.StatusCode).To(Equal(fiber.StatusCreated))
			Expect(body).To(MatchJSON(expected))
		})
	})

	Describe("GET /ticket/{id}", func() {
		var (
			id  string
			url string
		)

		BeforeEach(func() {
			id = strconv.Itoa(gofakeit.Number(1, 100))
			url = fmt.Sprintf(getTicketURL, id)
		})

		It("should return get ticket error when request is not validated", func() {
			url = fmt.Sprintf(getTicketURL, "a")

			prepareRequest(fiber.MethodGet, url, nil)

			actual, body := sendTestRequest(req)

			assertValidationFailed(actual, body)
		})
		It("should return get ticket error when ticket not found", func() {
			prepareRequest(fiber.MethodGet, url, nil)

			actual, body := sendTestRequest(req)

			Expect(actual.StatusCode).To(Equal(fiber.StatusBadRequest))
			assertErrorResponseByCodeWithMessage(body, "ticket_not_found", "Ticket not found")
		})

		It("should return get ticket response", func() {
			name := gofakeit.Name()
			desc := gofakeit.SentenceSimple()
			allocation := gofakeit.IntRange(1, 100)
			rows, err := mysqlInstance.Database().Exec("insert into ticket (name, description, quantity) values (?, ?, ?)", name, desc, allocation)

			Expect(err).To(BeNil())

			id, err := rows.LastInsertId()
			Expect(err).To(BeNil())

			expected := fmt.Sprintf(`{
			  "data": {
				"id": %d,
				"name": "%s",
				"desc": "%s",
				"allocation": %d
			  }
			}`, id, name, desc, allocation)

			url = fmt.Sprintf(getTicketURL, strconv.Itoa(int(id)))

			prepareRequest(fiber.MethodGet, url, nil)

			actual, body := sendTestRequest(req)

			Expect(actual.StatusCode).To(Equal(fiber.StatusOK))
			Expect(body).To(MatchJSON(expected))
		})
	})

	Describe("GET /v1/ticket_options/{id}/purchase", func() {
		var (
			id    string
			url   string
			ptReq request.PurchaseTicketRequest
		)

		BeforeEach(func() {
			id = strconv.Itoa(gofakeit.Number(1, 100))
			url = fmt.Sprintf(purcahseTicketURL, id)
			ptReq = request.PurchaseTicketRequest{
				ID:       id,
				Quantity: gofakeit.IntRange(3, 100),
				UserID:   gofakeit.UUID(),
			}
		})

		It("should return purchase ticket error when ticket not found", func() {
			prepareRequest(fiber.MethodPost, url, prepareBody(ptReq))

			actual, body := sendTestRequest(req)

			Expect(actual.StatusCode).To(Equal(fiber.StatusBadRequest))
			assertErrorResponseByCodeWithMessage(body, "ticket_not_found", "Ticket not found")
		})

		It("should return purchase ticket error when ticket not enough", func() {
			name := gofakeit.Name()
			desc := gofakeit.SentenceSimple()
			allocation := 1
			rows, err := mysqlInstance.Database().Exec("insert into ticket (name, description, quantity) values (?, ?, ?)", name, desc, allocation)
			Expect(err).To(BeNil())

			id, err := rows.LastInsertId()
			Expect(err).To(BeNil())

			url = fmt.Sprintf(purcahseTicketURL, strconv.Itoa(int(id)))

			prepareRequest(fiber.MethodPost, url, prepareBody(ptReq))

			actual, body := sendTestRequest(req)

			Expect(actual.StatusCode).To(Equal(fiber.StatusBadRequest))
			assertErrorResponseByCodeWithMessage(body, "ticket_not_found", "Ticket not found")
		})

		It("should return empty payload response", func() {
			name := gofakeit.Name()
			desc := gofakeit.SentenceSimple()
			allocation := gofakeit.IntRange(1, 100)
			rows, err := mysqlInstance.Database().Exec("insert into ticket (name, description, quantity) values (?, ?, ?)", name, desc, allocation)

			Expect(err).To(BeNil())

			id, err := rows.LastInsertId()
			ptReq.ID = strconv.Itoa(int(id))
			Expect(err).To(BeNil())

			url = fmt.Sprintf(purcahseTicketURL, strconv.Itoa(int(id)))

			prepareRequest(fiber.MethodPost, url, prepareBody(ptReq))

			actual, body := sendTestRequest(req)

			Expect(actual.StatusCode).To(Equal(fiber.StatusOK))
			Expect(body).To(MatchJSON(emptyPayload))
		})
	})
})
