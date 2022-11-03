package api

import (
	"banksystem/util"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var CurrencyValidator validator.Func = func(field validator.FieldLevel) bool {
	v, ok := field.Field().Interface().(string)

	if !ok {
		fmt.Printf("Can't parse value %s", v)
		return false
	}

	return util.IsValidCurrency(v)
}
