package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/thaian1234/simplebank/utils"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return utils.IsSupportedCurrencies(currency)
	}
	return false
}
