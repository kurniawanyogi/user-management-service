package common

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"strings"
	"user-management-service/model"
)

var ErrValidator = map[string]string{}

func NewValidator() *validator.Validate {
	validate := validator.New()
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	enTranslations.RegisterDefaultTranslations(validate, trans)

	return validate
}

func FormatValidationError(err error) (validationResponses []model.ValidationResponse) {

	if fieldErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				validationResponses = append(validationResponses, model.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is a required field", err.Field()),
				})
			case "len":
				validationResponses = append(validationResponses, model.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s must be a %s length", err.Field(), err.Param()),
				})
			case "min":
				validationResponses = append(validationResponses, model.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s must be a minimum of %s in length", err.Field(), err.Param()),
				})
			case "max":
				validationResponses = append(validationResponses, model.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s must be a maximum of %s in length", err.Field(), err.Param()),
				})
			case "url":
				validationResponses = append(validationResponses, model.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s must be a valid URL", err.Field()),
				})
			case "oneof":
				validationResponses = append(validationResponses, model.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s must be an oneof [%s]", err.Field(), err.Param()),
				})
			case "required_if":
				params := strings.Split(err.Param(), " ")
				formattedParams := params[0]
				for i, param := range params {
					if i > 0 {
						if i%2 != 0 {
							formattedParams += fmt.Sprintf(" is %s", param)
						} else {
							formattedParams += fmt.Sprintf(" and %s", param)
						}
					}
				}
				validationResponses = append(validationResponses, model.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is a required if %s", err.Field(), formattedParams),
				})
			case "required_unless":
				paramString := err.Param()
				formattedParams := strings.Replace(paramString, " ", " is not ", -1)
				validationResponses = append(validationResponses, model.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is a required if %s", err.Field(), formattedParams),
				})
			case "required_without":
				validationResponses = append(validationResponses, model.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is a required if %s is empty", err.Field(), err.Param()),
				})
			case "required_without_all":
				validationResponses = append(validationResponses, model.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is a required if %s are empty", err.Field(), err.Param()),
				})
			case "required_with":
				validationResponses = append(validationResponses, model.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is a required if %s is not empty", err.Field(), err.Param()),
				})
			case "excluded_with":
				validationResponses = append(validationResponses, model.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is a exclude if %s is empty", err.Field(), err.Param()),
				})
			case "ltecsfield":
				validationResponses = append(validationResponses, model.ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is less than to another %s field", err.Field(), err.Param()),
				})
			default:
				errValidator, ok := ErrValidator[err.Tag()]
				if ok {
					count := strings.Count(errValidator, "%s")
					if count == 1 {
						validationResponses = append(validationResponses, model.ValidationResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(errValidator, err.Field()),
						})
					} else {
						validationResponses = append(validationResponses, model.ValidationResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(errValidator, err.Field(), err.Param()),
						})
					}
				} else {
					validationResponses = append(validationResponses, model.ValidationResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf("something wrong on %s; %s", err.Field(), err.Tag()),
					})
				}
			}
		}
	}
	return validationResponses
}
