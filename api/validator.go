package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/radugaf/simplebank/tools"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return tools.IsSupportedCurrency(currency)
	}
	return false
}
