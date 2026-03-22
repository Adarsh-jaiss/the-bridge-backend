package validator

import (
	"net/url"
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
	v.RegisterAlias("safe_alphabets", "alphaspace")
	v.RegisterAlias("safe_url", "https_url")
	v.RegisterAlias("safe_alphabets_with_numbers", "alphanum")
	v.RegisterAlias("safe_int", "min=1")
}

func registerCustomValidators(v *validator.Validate) {
	htmlTagRegex := regexp.MustCompile(`<[^>]+>`)
	v.RegisterValidation("no_html", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		return !htmlTagRegex.MatchString(value)
	})

	v.RegisterValidation("https_url", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()

		if value == "" {
			return true
		}

		parsed, err := url.ParseRequestURI(value)
		if err != nil {
			return false
		}

		return parsed.Scheme == "https" && parsed.Host != ""
	})
}
