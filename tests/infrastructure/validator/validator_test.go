package validator_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/validator"
)

func TestValidator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Validator Suite")
}

var _ = Describe("Validator", func() {
	var v validator.Validator

	BeforeEach(func() {
		v = validator.NewValidator()
	})

	Describe("NewValidator", func() {
		It("should create validator instance", func() {
			Expect(v).ToNot(BeNil())
		})
	})

	Describe("Validate", func() {
		It("should validate struct with valid data", func() {
			type TestStruct struct {
				Name  string `validate:"required"`
				Email string `validate:"required,email"`
			}

			testData := TestStruct{
				Name:  "John Doe",
				Email: "john@example.com",
			}

			err := v.Validate(testData)
			Expect(err).To(BeNil())
		})

		It("should return error for invalid data", func() {
			type TestStruct struct {
				Name  string `validate:"required"`
				Email string `validate:"required,email"`
			}

			testData := TestStruct{
				Name:  "", // Required field is empty
				Email: "invalid-email", // Invalid email format
			}

			err := v.Validate(testData)
			Expect(err).ToNot(BeNil())
			// Validate returns raw validator errors, not formatted ones
			Expect(err.Error()).To(ContainSubstring("failed on the 'required' tag"))
		})
	})

	Describe("ValidateStruct", func() {
		It("should validate struct and return formatted errors", func() {
			type TestStruct struct {
				Name  string `validate:"required"`
				Email string `validate:"required,email"`
			}

			testData := TestStruct{
				Name:  "",
				Email: "invalid-email",
			}

			err := v.ValidateStruct(testData)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(ContainSubstring("Field 'Name' failed validation: required"))
			Expect(err.Error()).To(ContainSubstring("Field 'Email' failed validation: email"))
		})
	})

	Describe("ValidateField", func() {
		It("should validate single field", func() {
			err := v.ValidateField("john@example.com", "email")
			Expect(err).To(BeNil())
		})

		It("should return error for invalid field", func() {
			err := v.ValidateField("invalid-email", "email")
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Validation Constants", func() {
		It("should have validation constants", func() {
			Expect(validator.Required).To(Equal("required"))
			Expect(validator.Email).To(Equal("email"))
			Expect(validator.Min).To(Equal("min"))
			Expect(validator.Max).To(Equal("max"))
			Expect(validator.URL).To(Equal("url"))
			Expect(validator.UUID).To(Equal("uuid"))
		})
	})

	Describe("ValidationError", func() {
		It("should create validation error", func() {
			validationErr := validator.ValidationError{
				Field:   "email",
				Tag:     "email",
				Value:   "invalid-email",
				Message: "Invalid email format",
			}

			Expect(validationErr.Field).To(Equal("email"))
			Expect(validationErr.Tag).To(Equal("email"))
			Expect(validationErr.Value).To(Equal("invalid-email"))
			Expect(validationErr.Message).To(Equal("Invalid email format"))
		})
	})

	Describe("ValidationErrors", func() {
		It("should handle multiple validation errors", func() {
			errors := validator.ValidationErrors{
				{
					Field:   "name",
					Tag:     "required",
					Message: "Name is required",
				},
				{
					Field:   "email",
					Tag:     "email",
					Message: "Invalid email format",
				},
			}

			Expect(errors.Error()).To(ContainSubstring("Name is required"))
			Expect(errors.Error()).To(ContainSubstring("Invalid email format"))

			errorMap := errors.ToMap()
			Expect(errorMap["count"]).To(Equal(2))
			Expect(errorMap["errors"]).To(Equal(errors))
		})
	})
}) 