package request

import (
	"strconv"

	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket/model"
)

type CreateTicketRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"desc" validate:"required"`
	Allocation  int    `json:"allocation" validate:"required,number,gt=0"`
}

type GetTicketRequest struct {
	ID string `json:"id" validate:"required,number"`
}

type PurchaseTicketRequest struct {
	ID       string `json:"id" validate:"required,number" swaggerignore:"true"`
	Quantity int    `json:"quantity" validate:"required,number,gt=0"`
	UserID   string `json:"user_id" validate:"required,uuid"`
}

func (ctr CreateTicketRequest) ToPayload() model.CreateTicketRequest {
	return model.CreateTicketRequest{
		Name:        ctr.Name,
		Description: ctr.Description,
		Allocation:  ctr.Allocation,
	}
}

func (gtr GetTicketRequest) ToPayload() model.GetTicketRequest {
	id, _ := strconv.ParseInt(gtr.ID, 10, 64) // nolint: gomnd

	return model.GetTicketRequest{
		ID: id,
	}
}

func (ptr PurchaseTicketRequest) ToPayload() model.PurchaseTicketRequest {
	id, _ := strconv.ParseInt(ptr.ID, 10, 64) // nolint: gomnd

	return model.PurchaseTicketRequest{
		ID:       id,
		Quantity: ptr.Quantity,
		UserID:   ptr.UserID,
	}
}
