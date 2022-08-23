package utils

import (
	"context"

	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/validation"
)

func ValidateWithContext(ctx context.Context, i interface{}) map[string]string {
	v := ctx.Value(ValidatorKey).(validation.IValidator)
	if errs := v.Validate(i); len(errs) > 0 {
		return errs
	}

	return nil
}
