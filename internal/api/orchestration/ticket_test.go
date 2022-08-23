package orchestration_test

import (
	"context"
	"fmt"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/dto/request"
	. "github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/orchestration"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket/mocks"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket/model"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/constants"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/utils"
)

var _ = Describe("TicketOrchestration", func() {
	var (
		ctx               context.Context
		mockCtrl          *gomock.Controller
		ticketServiceMock *mocks.MockITicketService
		orchestrator      ITicketOrchestrator
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockCtrl = gomock.NewController(GinkgoT())
		ticketServiceMock = mocks.NewMockITicketService(mockCtrl)
		orchestrator = NewTicketOrchestrator(ticketServiceMock)
	})

	Describe("Create", func() {
		var (
			ticket request.CreateTicketRequest
			res    *model.CreateTicketResource
		)

		BeforeEach(func() {
			ticket = request.CreateTicketRequest{
				Name:        gofakeit.Name(),
				Description: gofakeit.SentenceSimple(),
				Allocation:  gofakeit.IntRange(1, 100),
			}
			res = &model.CreateTicketResource{
				ID:          gofakeit.Int64(),
				Name:        ticket.Name,
				Description: ticket.Description,
				Allocation:  ticket.Allocation,
			}
		})

		It("should return error when ticket service returns error", func() {
			ticketServiceMock.EXPECT().Create(ctx, ticket.ToPayload()).Return(nil, errService)

			actual, err := orchestrator.Create(ctx, ticket)

			Expect(err).To(HaveOccurred())
			Expect(actual).To(BeNil())
		})

		It("should return ticket", func() {
			ticketServiceMock.EXPECT().Create(ctx, ticket.ToPayload()).Return(res, nil)

			actual, err := orchestrator.Create(ctx, ticket)

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(res))
		})
	})

	Describe("Get", func() {
		var (
			ticket request.GetTicketRequest
			res    *model.GetTicketResource
		)

		BeforeEach(func() {
			ticket = request.GetTicketRequest{
				ID: strconv.Itoa(gofakeit.IntRange(1, 100)),
			}
			res = &model.GetTicketResource{
				ID:          gofakeit.Int64(),
				Name:        gofakeit.Name(),
				Description: gofakeit.SentenceSimple(),
				Allocation:  gofakeit.IntRange(1, 100),
			}
		})

		It("should return error when ticket service returns error", func() {
			ticketServiceMock.EXPECT().Get(ctx, ticket.ToPayload()).Return(nil, errService)

			actual, err := orchestrator.Get(ctx, ticket)

			Expect(err).To(HaveOccurred())
			Expect(actual).To(BeNil())
		})

		It("should return ticket not found error", func() {
			ticketServiceMock.EXPECT().Get(ctx, ticket.ToPayload()).Return(nil, fmt.Errorf(constants.TicketNotFound))

			expected := utils.ErrorBag{
				Code:    constants.TicketNotFound,
				Cause:   fmt.Errorf(constants.TicketNotFound),
				Message: "Ticket not found",
			}

			actual, err := orchestrator.Get(ctx, ticket)

			Expect(err).To(Equal(expected))
			Expect(actual).To(BeNil())
		})

		It("should return ticket", func() {
			ticketServiceMock.EXPECT().Get(ctx, ticket.ToPayload()).Return(res, nil)

			actual, err := orchestrator.Get(ctx, ticket)

			Expect(err).ToNot(HaveOccurred())
			Expect(actual).To(Equal(res))
		})
	})

	Describe("Purchase", func() {
		var (
			ticket request.PurchaseTicketRequest
		)

		BeforeEach(func() {
			ticket = request.PurchaseTicketRequest{
				UserID:   gofakeit.UUID(),
				Quantity: gofakeit.IntRange(1, 100),
				ID:       strconv.Itoa(gofakeit.IntRange(1, 100)),
			}
		})

		It("should return error when ticket service returns error", func() {
			ticketServiceMock.EXPECT().Purchase(ctx, ticket.ToPayload()).Return(errService)

			err := orchestrator.Purchase(ctx, ticket)

			Expect(err).To(HaveOccurred())
		})

		It("should return ticket not found error", func() {
			ticketServiceMock.EXPECT().Purchase(ctx, ticket.ToPayload()).Return(fmt.Errorf(constants.TicketNotFound))

			expected := utils.ErrorBag{
				Code:    constants.TicketNotFound,
				Cause:   fmt.Errorf(constants.TicketNotFound),
				Message: "Ticket not found",
			}

			err := orchestrator.Purchase(ctx, ticket)

			Expect(err).To(Equal(expected))
		})

		It("should return not enough ticket error", func() {
			ticketServiceMock.EXPECT().Purchase(ctx, ticket.ToPayload()).Return(fmt.Errorf(constants.NotEnoughTickets))

			expected := utils.ErrorBag{
				Code:    constants.NotEnoughTickets,
				Cause:   fmt.Errorf(constants.NotEnoughTickets),
				Message: "Not enough tickets",
			}

			err := orchestrator.Purchase(ctx, ticket)

			Expect(err).To(Equal(expected))
		})

		It("should return nil", func() {
			ticketServiceMock.EXPECT().Purchase(ctx, ticket.ToPayload()).Return(nil)

			err := orchestrator.Purchase(ctx, ticket)

			Expect(err).ToNot(HaveOccurred())
		})
	})
})
