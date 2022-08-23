package handler

import (
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/dto/request"
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/dto/resource"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/response"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/utils"
	"github.com/gofiber/fiber/v2"

	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/orchestration"
)

type ITicketHandler interface {
	Create(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Purchase(c *fiber.Ctx) error
}

type ticketHandler struct {
	ticketOrchestrator orchestration.ITicketOrchestrator
}

func NewTicketHandler(to orchestration.ITicketOrchestrator) ITicketHandler {
	return &ticketHandler{
		ticketOrchestrator: to,
	}
}

// Create godoc
// @Tags api.ticket
// @ID api.ticket.create
// @Description creates a ticket
// @Accept  json
// @Produce  json
// @Param ticket body request.CreateTicketRequest true "ticket creating requirements"
// @Success 200 {object} response.HTTPSuccessResponse{data=resource.CreateTicketResource}
// @Failure 400 {object} response.HTTPErrorResponse
// @Failure 422 {object} response.HTTPValidationErrorResponse
// @Router /v1/ticket_options [POST]
func (t *ticketHandler) Create(c *fiber.Ctx) error {
	var req request.CreateTicketRequest
	if reqErr := c.BodyParser(&req); reqErr != nil {
		return response.HandleParserErrorResponse(c)
	}

	ctx := utils.GetContextFromFiber(c)
	if vErr := utils.ValidateWithContext(ctx, req); vErr != nil {
		return response.HandleValidationErrorResponse(c, vErr)
	}

	res, err := t.ticketOrchestrator.Create(ctx, req)
	if err != nil {
		return response.HandleErrorResponse(ctx, c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(response.NewSuccessResponse(resource.NewCreateTicketResource(res)))
}

// Get godoc
// @Tags api.ticket
// @ID api.ticket.get
// @Description get a ticket
// @Accept  json
// @Produce  json
// @Param id path string true "ticket id"
// @Success 200 {object} response.HTTPSuccessResponse{data=resource.GetTicketResource}
// @Failure 400 {object} response.HTTPErrorResponse
// @Failure 422 {object} response.HTTPValidationErrorResponse
// @Router /v1/ticket/{id} [GET]
func (t *ticketHandler) Get(c *fiber.Ctx) error {
	req := request.GetTicketRequest{
		ID: c.Params("id"),
	}

	ctx := utils.GetContextFromFiber(c)
	if vErr := utils.ValidateWithContext(ctx, req); vErr != nil {
		return response.HandleValidationErrorResponse(c, vErr)
	}

	res, err := t.ticketOrchestrator.Get(ctx, req)
	if err != nil {
		return response.HandleErrorResponse(ctx, c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(resource.NewGetTicketResource(res)))
}

// Purchase godoc
// @Tags api.ticket
// @ID api.ticket.purchase
// @Description purchase a ticket
// @Accept  json
// @Produce  json
// @Param id path string true "ticket id"
// @Param ticket body request.PurchaseTicketRequest true "purchase ticket requirements"
// @Success 200 {object} object
// @Failure 400 {object} response.HTTPErrorResponse
// @Failure 422 {object} response.HTTPValidationErrorResponse
// @Router /v1/ticket_options/{id} [POST]
func (t *ticketHandler) Purchase(c *fiber.Ctx) error {
	req := request.PurchaseTicketRequest{
		ID: c.Params("id"),
	}

	if reqErr := c.BodyParser(&req); reqErr != nil {
		return response.HandleParserErrorResponse(c)
	}

	ctx := utils.GetContextFromFiber(c)
	if vErr := utils.ValidateWithContext(ctx, req); vErr != nil {
		return response.HandleValidationErrorResponse(c, vErr)
	}

	err := t.ticketOrchestrator.Purchase(ctx, req)
	if err != nil {
		return response.HandleErrorResponse(ctx, c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.NewSuccessResponse(fiber.Map{}))
}
