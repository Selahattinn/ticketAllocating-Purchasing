package resource

import "github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket/model"

type CreateTicketResource struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	Allocation  int    `json:"allocation"`
}

type GetTicketResource struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"desc"`
	Allocation  int    `json:"allocation"`
}

func NewCreateTicketResource(t *model.CreateTicketResource) *CreateTicketResource {
	return &CreateTicketResource{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Allocation:  t.Allocation,
	}
}

func NewGetTicketResource(t *model.GetTicketResource) *GetTicketResource {
	return &GetTicketResource{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Allocation:  t.Allocation,
	}
}
