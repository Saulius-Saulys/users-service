package validation

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func ValidateJSONBody(c *gin.Context, body any) (error, int) {
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
			return errors.New(errMsg), http.StatusUnprocessableEntity
		case errors.As(err, &unmarshalErr):
			errMsg += fmt.Sprintf("JSON field '%s' failed validation. Was: %v, expected %s.", strings.ToLower(unmarshalErr.Field), unmarshalErr.Value, unmarshalErr.Type)

			return errors.New(errMsg), http.StatusBadRequest
		default:
			return errors.New("unknown error"), http.StatusInternalServerError
		}
	}
	return nil, 0
}

func operatorDescription(shortOperator string) string {
	switch shortOperator {
	case "gte":
		return "greater than or equal to"
	case "lte":
		return "lower than or equal to"
	case "min":
		return "length greater or equal to"
	case "uri":
		return "to be a valid URI"
	}

	return shortOperator
}
