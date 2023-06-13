package app

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidateAlphaUnicodeWithSpace(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	val = strings.ReplaceAll(val, " ", "")

	err := validator.New().Var(val, "alphaunicode")

	return err == nil
}

func ValidateBeforeNow(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	t, err := time.Parse("2006-01-02", val)

	if err != nil {
		return false
	}

	return t.Before(time.Now())
}
