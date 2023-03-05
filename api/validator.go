package api

import (
	"github.com/crackz/simple-bank/util"
	"github.com/go-playground/validator/v10"
)

func validateCurrency(fieldLevel validator.FieldLevel) bool {

	if value, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(value)
	}

	return false
}
