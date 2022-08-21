package utils

// System errors
const (
	NotFoundErrCode   = "not_found"
	ValidationErrCode = "validation_failed"
	UnexpectedErrCode = "unexpected_error"
	BodyParserErrCode = "body_parser_failed"

	NotFoundMsg   = "Not found!"
	UnexpectedMsg = "An unexpected error has occurred."
	ValidationMsg = "The given data was invalid."
	BodyParserMsg = "The given values could not be parsed."
)

type ErrorBag struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"cause"`
}

func (e ErrorBag) GetMessage() string {
	return e.Message
}

func (e ErrorBag) GetCode() string {
	return e.Code
}

func (e ErrorBag) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return ""
}
