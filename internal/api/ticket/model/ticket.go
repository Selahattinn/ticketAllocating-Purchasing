package model

type CreateTicketRequest struct {
	Name        string `json:"name"`
	Description string `json:"desc"`
	Allocation  int    `json:"allocation"`
}

type GetTicketRequest struct {
	ID int64 `json:"id"`
}

type PurchaseTicketRequest struct {
	ID       int64  `json:"id"`
	Quantity int    `json:"quantity"`
	UserID   string `json:"user_id"`
}

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

func (gtr *GetTicketResource) IsEmpty() bool {
	return gtr.ID == 0
}
