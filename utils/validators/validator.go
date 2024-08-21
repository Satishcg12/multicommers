package validators

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		errorMessages := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			jsonTag := cv.getJSONTag(i, err.StructField())
			if jsonTag != "" {
				errorMessages[jsonTag] = cv.getErrorMsg(err)
			} else {
				errorMessages[err.Field()] = cv.getErrorMsg(err)
			}
		}
		return echo.NewHTTPError(http.StatusBadRequest, errorMessages)
	}
	return nil
}

func (cv *CustomValidator) getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters"
	case "email":
		return "Invalid email format"
	case "eqfield":
		return fe.Field() + " must be equal to " + fe.Param()
	case "fullname":
		return "Full name should contain first name and last name (middle name optional)"
	case "password":
		return "Password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character"
	default:
		return fe.Field() + " is invalid"
	}
}

func (cv *CustomValidator) getJSONTag(model interface{}, fieldName string) string {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	field, found := modelType.FieldByName(fieldName)
	if !found {
		return ""
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName
	}
	jsonTag = strings.Split(jsonTag, ",")[0]
	return jsonTag
}

func NewValidator() *CustomValidator {
	v := validator.New()
	v.RegisterValidation("fullname", validateFullName) // Register the custom validator
	v.RegisterValidation("password", validatePassword)
	return &CustomValidator{Validator: v}
}

var fullNameRegex = regexp.MustCompile(`^[a-zA-Z]+([ '-][a-zA-Z]+){1,2}$`)

func validateFullName(fl validator.FieldLevel) bool {
	return fullNameRegex.MatchString(fl.Field().String())
}

// Password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// Check length
	if len(password) < 8 {
		return false
	}

	// Check for at least one lowercase letter
	lowercase := regexp.MustCompile(`[a-z]`)
	if !lowercase.MatchString(password) {
		return false
	}

	// Check for at least one uppercase letter
	uppercase := regexp.MustCompile(`[A-Z]`)
	if !uppercase.MatchString(password) {
		return false
	}

	// Check for at least one digit
	digit := regexp.MustCompile(`\d`)
	return digit.MatchString(password)
}
