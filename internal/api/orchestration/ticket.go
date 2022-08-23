package orchestration

import (
	"context"

	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/dto/request"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket/model"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/constants"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/utils"
)

type ITicketOrchestrator interface {
	Create(ctx context.Context, req request.CreateTicketRequest) (*model.CreateTicketResource, error)
	Get(ctx context.Context, req request.GetTicketRequest) (*model.GetTicketResource, error)
	Purchase(ctx context.Context, req request.PurchaseTicketRequest) error
}

type ticketOrchestrator struct {
	ticketService ticket.ITicketService
}

func NewTicketOrchestrator(to ticket.ITicketService) ITicketOrchestrator {
	return &ticketOrchestrator{

		ticketService: to,
	}
}

func (t *ticketOrchestrator) Create(
	ctx context.Context, req request.CreateTicketRequest,
) (*model.CreateTicketResource, error) {
	ti, err := t.ticketService.Create(ctx, req.ToPayload())
	if err != nil {
		return nil, utils.ErrorBag{Code: constants.ProcessingErrCode, Cause: err}
	}

	return ti, nil
}

func (t *ticketOrchestrator) Get(ctx context.Context, req request.GetTicketRequest) (*model.GetTicketResource, error) {
	ti, err := t.ticketService.Get(ctx, req.ToPayload())
	if err != nil {
		if err.Error() == constants.TicketNotFound {
			return nil, utils.ErrorBag{Code: constants.TicketNotFound, Cause: err, Message: "Ticket not found"}
		}

		return nil, utils.ErrorBag{Code: constants.ProcessingErrCode, Cause: err}
	}

	return ti, nil
}

func (t *ticketOrchestrator) Purchase(ctx context.Context, req request.PurchaseTicketRequest) error {
	if err := t.ticketService.Purchase(ctx, req.ToPayload()); err != nil {
		if err.Error() == constants.TicketNotFound {
			return utils.ErrorBag{Code: constants.TicketNotFound, Cause: err, Message: "Ticket not found"}
		}

		if err.Error() == constants.NotEnoughTickets {
			return utils.ErrorBag{Code: constants.NotEnoughTickets, Cause: err, Message: "Not enough tickets"}
		}

		return utils.ErrorBag{Code: constants.ProcessingErrCode, Cause: err}
	}

	return nil
}
