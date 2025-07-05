package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(interface{}) error
	ValidateStruct(interface{}) error
	ValidateField(interface{}, string) error
}

type validatorImpl struct {
	validate *validator.Validate
}

func NewValidator() Validator {
	v := validator.New()
	
	// v.RegisterValidation("custom_rule", customValidationFunc)
	
	return &validatorImpl{
		validate: v,
	}
}

func (v *validatorImpl) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func (v *validatorImpl) ValidateStruct(i interface{}) error {
	err := v.validate.Struct(i)
	if err != nil {
		return v.formatValidationErrors(err)
	}
	return nil
}

func (v *validatorImpl) ValidateField(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}

func (v *validatorImpl) formatValidationErrors(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()
			
			message := fmt.Sprintf("Field '%s' failed validation: %s", field, tag)
			if param != "" {
				message += fmt.Sprintf(" (value: %s)", param)
			}
			
			errorMessages = append(errorMessages, message)
		}
		
		return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, "; "))
	}
	
	return err
}

const (
	Required = "required"
	Email    = "email"
	Min      = "min"
	Max      = "max"
	URL      = "url"
	UUID     = "uuid"
)

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var messages []string
	for _, err := range v {
		messages = append(messages, err.Message)
	}
	return strings.Join(messages, "; ")
}

func (v ValidationErrors) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"errors": v,
		"count":  len(v),
	}
}