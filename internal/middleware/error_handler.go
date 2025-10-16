package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only handle errors if the response has not been written yet
		if c.Writer.Status() == http.StatusOK {
			if len(c.Errors) > 0 {
				err := c.Errors.Last().Err
				handleError(c, err)
			}
		}
	}
}

func handleError(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		handleValidationError(c, validationErrors)
		return
	}

	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: "Internal server error",
	})
}

func handleValidationError(c *gin.Context, validationErrors validator.ValidationErrors) {
	errs := make([]ValidationError, len(validationErrors))
	for i, fieldError := range validationErrors {
		errs[i] = ValidationError{
			Field:   fieldError.Field(),
			Message: getValidationErrorMessage(fieldError),
		}
	}

	c.JSON(http.StatusBadRequest, ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: "Validation failed",
		Errors:  errs,
	})
}

func getValidationErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Value must be greater than or equal to " + fieldError.Param()
	case "max":
		return "Value must be less than or equal to " + fieldError.Param()
	case "email":
		return "Invalid email format"
	case "url":
		return "Invalid URL format"
	case "e164":
		return "Invalid phone number format"
	default:
		return "Invalid value"
	}
}
