package ticket

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket/model"
)

type ITicketService interface {
	Create(ctx context.Context, ticket model.CreateTicketRequest) (*model.CreateTicketResource, error)
	Get(ctx context.Context, ticket model.GetTicketRequest) (*model.GetTicketResource, error)
	Purchase(ctx context.Context, ticket model.PurchaseTicketRequest) error
}

type ticketService struct {
	logger     *logrus.Logger
	repository ITicketRepository
}

func NewTicketService(l *logrus.Logger, tr ITicketRepository) ITicketService {
	return &ticketService{
		logger:     l,
		repository: tr,
	}
}

func (t *ticketService) Create(
	ctx context.Context, ticket model.CreateTicketRequest,
) (*model.CreateTicketResource, error) {
	id, err := t.repository.Create(ctx, ticket)
	if err != nil {
		t.logger.WithContext(ctx).WithError(err).Error("create ticket: repository err")
		return nil, err
	}

	return &model.CreateTicketResource{
		ID:          id,
		Name:        ticket.Name,
		Description: ticket.Description,
		Allocation:  ticket.Allocation,
	}, nil
}

func (t *ticketService) Get(ctx context.Context, ticket model.GetTicketRequest) (*model.GetTicketResource, error) {
	ticketRes, err := t.repository.Get(ctx, ticket)
	if err != nil {
		t.logger.WithContext(ctx).WithError(err).Error("get ticket: repository err")
		return nil, err
	}

	return ticketRes, nil
}

func (t *ticketService) Purchase(ctx context.Context, ticket model.PurchaseTicketRequest) error {
	err := t.repository.Purchase(ctx, ticket)
	if err != nil {
		t.logger.WithContext(ctx).WithError(err).Error("get ticket: repository err")
		return err
	}

	return nil
}
