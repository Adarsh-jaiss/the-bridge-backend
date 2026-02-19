package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func RegisterOnValidator(v *validator.Validate) {
	registerCustomValidators(v)
	registerAliases(v)
}

func registerAliases(v *validator.Validate) {
	v.RegisterAlias("safe_email", "email")
	v.RegisterAlias("safe_otp", "max=6,numeric")
}

func registerCustomValidators(v *validator.Validate) {
	htmlTagRegex := regexp.MustCompile(`<[^>]+>`)
	v.RegisterValidation("no_html", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		return !htmlTagRegex.MatchString(value)
	})
}
