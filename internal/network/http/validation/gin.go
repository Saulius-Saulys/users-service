package validation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

func ValidateJSONBody(c *gin.Context, body any) (int, error) {
	var errMsg string

	if err := c.ShouldBindJSON(body); err != nil {
		var validationErr validator.ValidationErrors
		var unmarshalErr *json.UnmarshalTypeError
		switch {
		case errors.As(err, &validationErr):
			for _, fieldErr := range validationErr {
				description := operatorDescription(fieldErr.Tag())
				errMsg += fmt.Sprintf("JSON field '%s' failed validation. Was: %v, expected %s %s.", strings.ToLower(fieldErr.Field()), fieldErr.Value(), description, fieldErr.Param())
			}
			return http.StatusUnprocessableEntity, errors.New(errMsg)
		case errors.As(err, &unmarshalErr):
			errMsg += fmt.Sprintf("JSON field '%s' failed validation. Was: %v, expected %s.", strings.ToLower(unmarshalErr.Field), unmarshalErr.Value, unmarshalErr.Type)

			return http.StatusBadRequest, errors.New(errMsg)
		default:
			return http.StatusInternalServerError, errors.New("unknown error")
		}
	}
	return 0, nil
}

func operatorDescription(shortOperator string) string {
	switch shortOperator {
	case "min":
		return "greater than"
	case "max":
		return "lower than"
	case "email":
		return "to be a valid email"
	}

	return shortOperator
}
