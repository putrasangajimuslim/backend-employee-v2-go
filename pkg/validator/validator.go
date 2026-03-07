// pkg/validator/validator.go
package validator

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) map[string][]string {
	if validate == nil {
		Init()
	}

	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string][]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()
		param := err.Param()

		var message string
		switch tag {
		case "required":
			message = fmt.Sprintf("%s is required", field)
		case "email":
			message = fmt.Sprintf("%s must be a valid email", field)
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters", field, param)
		case "max":
			message = fmt.Sprintf("%s must be at most %s characters", field, param)
		default:
			message = fmt.Sprintf("%s is invalid", field)
		}

		errors[field] = append(errors[field], message)
	}

	return errors
}

func GetValidator() *validator.Validate {
	if validate == nil {
		Init()
	}
	return validate
}

func AddValidation(tag string, fn validator.Func) error {
	if validate == nil {
		Init()
	}
	return validate.RegisterValidation(tag, fn)
}

func StructFieldName(s interface{}, field string) string {
	t := reflect.TypeOf(s)
	f, found := t.FieldByName(field)
	if !found {
		return field
	}
	return f.Name
}
