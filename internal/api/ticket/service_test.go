package ticket_test

import (
	"context"
	"errors"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	. "github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket/mocks"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket/model"
)

var (
	errmysql = errors.New("mysql err")
)

var _ = Describe("TicketService", func() {
	var (
		ctx            context.Context
		logger         *logrus.Logger
		mockCtrl       *gomock.Controller
		repositoryMock *mocks.MockITicketRepository
		service        ITicketService
	)

	BeforeEach(func() {
		ctx = context.Background()
		logger, _ = test.NewNullLogger()
		mockCtrl = gomock.NewController(GinkgoT())
		repositoryMock = mocks.NewMockITicketRepository(mockCtrl)
		service = NewTicketService(logger, repositoryMock)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Create", func() {
		var (
			id     int64
			ticket model.CreateTicketRequest
		)

		BeforeEach(func() {
			ticket = model.CreateTicketRequest{
				Name:        gofakeit.Name(),
				Description: gofakeit.SentenceSimple(),
				Allocation:  gofakeit.IntRange(1, 100),
			}
			id = gofakeit.Int64()
		})

		It("should return error when repository returns error", func() {
			repositoryMock.EXPECT().Create(ctx, ticket).Return(int64(-1), errmysql)

			_, err := service.Create(ctx, ticket)

			Expect(err).To(HaveOccurred())
		})

		It("should return ticket", func() {
			repositoryMock.EXPECT().Create(ctx, ticket).Return(id, nil)

			expected := &model.CreateTicketResource{
				ID:          id,
				Name:        ticket.Name,
				Description: ticket.Description,
				Allocation:  ticket.Allocation,
			}

			res, err := service.Create(ctx, ticket)

			Expect(err).To(BeNil())
			Expect(res).To(Equal(expected))
		})
	})

	Describe("Get", func() {
		var (
			id        int64
			ticket    model.GetTicketRequest
			ticketRes *model.GetTicketResource
		)

		BeforeEach(func() {
			id = gofakeit.Int64()
			ticket = model.GetTicketRequest{
				ID: id,
			}
			ticketRes = &model.GetTicketResource{
				ID:          id,
				Name:        gofakeit.Name(),
				Description: gofakeit.SentenceSimple(),
				Allocation:  gofakeit.IntRange(1, 100),
			}
		})

		It("should return error when repository returns error", func() {
			repositoryMock.EXPECT().Get(ctx, ticket).Return(ticketRes, errmysql)

			res, err := service.Get(ctx, ticket)

			Expect(err).To(HaveOccurred())
			Expect(res).To(BeNil())
		})

		It("should return ticket", func() {
			repositoryMock.EXPECT().Get(ctx, ticket).Return(ticketRes, nil)

			res, err := service.Get(ctx, ticket)

			Expect(err).To(BeNil())
			Expect(res).To(Equal(ticketRes))
		})
	})

	Describe("Purchase", func() {
		var (
			id     int64
			ticket model.PurchaseTicketRequest
		)

		BeforeEach(func() {
			id = gofakeit.Int64()
			ticket = model.PurchaseTicketRequest{
				ID:       id,
				Quantity: gofakeit.IntRange(1, 100),
				UserID:   gofakeit.UUID(),
			}
		})

		It("should return error when repository returns error", func() {
			repositoryMock.EXPECT().Purchase(ctx, ticket).Return(errmysql)

			err := service.Purchase(ctx, ticket)

			Expect(err).To(HaveOccurred())
		})

		It("should return nil", func() {
			repositoryMock.EXPECT().Purchase(ctx, ticket).Return(nil)

			err := service.Purchase(ctx, ticket)

			Expect(err).To(BeNil())
		})
	})
})
