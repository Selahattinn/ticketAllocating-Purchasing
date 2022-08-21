package response

import (
	"context"
	"github.com/gofiber/fiber/v2"

	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/utils"
)

type ErrorAttribute struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

type ErrorSchema struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type HTTPSuccessResponse struct {
	Data interface{} `json:"data"`
}

type HTTPErrorResponse struct {
	Error ErrorSchema `json:"error"`
}

type HTTPValidationErrorResponse struct {
	Error      ErrorSchema      `json:"error"`
	Attributes []ErrorAttribute `json:"attributes"`
}

func NewSuccessResponse(data interface{}) HTTPSuccessResponse {
	return HTTPSuccessResponse{
		Data: data,
	}
}

func NewErrorResponse(ctx context.Context, err error, msg ...string) HTTPErrorResponse {
	schema := ErrorSchema{
		Code:    utils.UnexpectedErrCode,
		Message: utils.UnexpectedMsg,
	}

	if errorBag, ok := err.(utils.ErrorBag); ok {
		schema.Code = errorBag.GetCode()

		switch {
		case errorBag.GetMessage() != "":
			schema.Message = errorBag.GetMessage()
		case len(msg) > 0:
			schema.Message = msg[0]
		default:
			schema.Message = utils.UnexpectedMsg
		}
	}

	return HTTPErrorResponse{Error: schema}
}

func NewBodyParserErrorResponse() HTTPErrorResponse {
	return HTTPErrorResponse{
		Error: ErrorSchema{
			Code:    utils.BodyParserErrCode,
			Message: utils.BodyParserMsg,
		},
	}
}

func NewValidationErrorResponse(errors map[string]string) HTTPValidationErrorResponse {
	var attrs []ErrorAttribute
	for k, v := range errors {
		attrs = append(attrs, ErrorAttribute{
			Name:    k,
			Message: v,
		})
	}

	return HTTPValidationErrorResponse{
		Error: ErrorSchema{
			Code:    utils.ValidationErrCode,
			Message: utils.ValidationMsg,
		},
		Attributes: attrs,
	}
}

func HandleErrorResponse(ctx context.Context, c *fiber.Ctx, err error, msg ...string) error {
	return c.Status(fiber.StatusBadRequest).JSON(NewErrorResponse(ctx, err, msg...))
}

func HandleValidationErrorResponse(c *fiber.Ctx, errors map[string]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(NewValidationErrorResponse(errors))
}

func HandleParserErrorResponse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(NewBodyParserErrorResponse())
}
