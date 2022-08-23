package ticket_test

import (
	"context"
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"

	. "github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket/model"
)

var _ = Describe("TicketRepository", func() {
	var (
		ctx        context.Context
		cancelCtx  context.CancelFunc
		repository ITicketRepository
	)

	BeforeEach(func() {
		ctx = context.Background()
		repository = NewTicketRepository(mysqlInstance)
	})

	AfterEach(func() {
		_, _ = mysqlInstance.Database().Exec("DELETE FROM ticket")
		_, _ = mysqlInstance.Database().Exec("DELETE FROM purchase")

	})

	cancelMysql := func() {
		ctx, cancelCtx = context.WithCancel(ctx)
		cancelCtx()
	}

	Describe("Get", func() {
		var (
			ticket    model.GetTicketRequest
			ticketRes *model.GetTicketResource
		)

		BeforeEach(func() {
			ticket = model.GetTicketRequest{
				ID: gofakeit.Int64(),
			}
			ticketRes = &model.GetTicketResource{
				ID:          ticket.ID,
				Name:        gofakeit.Name(),
				Description: gofakeit.SentenceSimple(),
				Allocation:  gofakeit.IntRange(1, 100),
			}
		})

		It("should return mysql error", func() {
			cancelMysql()

			actual, err := repository.Get(ctx, ticket)

			Expect(err).To(HaveOccurred())
			Expect(actual).To(BeNil())
		})

		It("should return ticket_not_found error", func() {
			expected := fmt.Errorf("ticket_not_found")

			actual, err := repository.Get(ctx, ticket)

			Expect(actual).To(BeNil())
			Expect(err).To(Equal(expected))
		})

		It("should return ticket", func() {
			res, _ := mysqlInstance.Database().Exec(
				"INSERT INTO ticket (name, description, quantity) VALUES (?, ?, ?)",
				ticketRes.Name, ticketRes.Description, ticketRes.Allocation)
			id, _ := res.LastInsertId()
			ticket.ID = id
			ticketRes.ID = id

			actual, err := repository.Get(ctx, ticket)

			Expect(err).To(BeNil())
			Expect(actual).To(Equal(ticketRes))
		})
	})

	Describe("Create", func() {
		var (
			ticket    model.CreateTicketRequest
			ticketRes model.CreateTicketResource
		)

		BeforeEach(func() {
			ticket = model.CreateTicketRequest{
				Name:        gofakeit.Name(),
				Description: gofakeit.SentenceSimple(),
				Allocation:  gofakeit.IntRange(1, 100),
			}
			ticketRes = model.CreateTicketResource{}
		})

		It("should return mysql error", func() {
			cancelMysql()

			actual, err := repository.Create(ctx, ticket)

			Expect(err).To(HaveOccurred())
			Expect(actual).To(Equal(int64(-1)))
		})

		It("should return ticket", func() {
			actual, err := repository.Create(ctx, ticket)

			errQ := mysqlInstance.Database().QueryRow("SELECT id, name, description, quantity FROM ticket WHERE id = ?", actual).
				Scan(&ticketRes.ID, &ticketRes.Name, &ticketRes.Description, &ticketRes.Allocation)
			if errQ != nil {
				panic(errQ)
			}

			Expect(err).To(BeNil())
			Expect(ticketRes).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"ID":          Equal(actual),
				"Name":        Equal(ticket.Name),
				"Description": Equal(ticket.Description),
				"Allocation":  Equal(ticket.Allocation),
			}))
		})
	})

	Describe("Purchase", func() {
		var (
			ticket model.PurchaseTicketRequest
		)

		BeforeEach(func() {
			ticket = model.PurchaseTicketRequest{
				ID:       1,
				UserID:   gofakeit.UUID(),
				Quantity: 2,
			}
		})

		It("should return mysql error", func() {
			cancelMysql()

			err := repository.Purchase(ctx, ticket)

			Expect(err).To(HaveOccurred())
		})

		It("should return ticket_not_found error", func() {
			expected := fmt.Errorf("ticket_not_found")

			err := repository.Purchase(ctx, ticket)

			Expect(err).To(Equal(expected))
		})

		It("should return not_enough_tickets error", func() {
			res, _ := mysqlInstance.Database().Exec("INSERT INTO ticket (name, description, quantity) VALUES (?, ?, ?)",
				gofakeit.Name(), gofakeit.SentenceSimple(), 1)
			id, _ := res.LastInsertId()
			ticket.ID = id

			expected := fmt.Errorf("not_enough_tickets")

			err := repository.Purchase(ctx, ticket)

			Expect(err).To(Equal(expected))
		})

		It("should return nil", func() {
			res, _ := mysqlInstance.Database().Exec("INSERT INTO ticket (name, description, quantity) VALUES (?, ?, ?)",
				gofakeit.Name(), gofakeit.SentenceSimple(), 3)
			id, _ := res.LastInsertId()
			ticket.ID = id

			err := repository.Purchase(ctx, ticket)

			var purchaseRes struct {
				Quantity int
				UserID   string
			}

			errQ := mysqlInstance.Database().QueryRow(
				"SELECT quantity, user_id FROM purchase WHERE user_id = ?",
				ticket.UserID).Scan(&purchaseRes.Quantity, &purchaseRes.UserID)
			if errQ != nil {
				panic(errQ)
			}

			Expect(err).To(BeNil())
			Expect(purchaseRes).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"UserID":   Equal(ticket.UserID),
				"Quantity": Equal(ticket.Quantity),
			}))
		})
	})
})
