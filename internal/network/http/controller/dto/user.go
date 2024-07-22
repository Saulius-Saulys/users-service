package dto

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type Country string

// expand the list of countries by adding more constants
const (
	UK  Country = "UK"
	LTU Country = "LTU"
	FR  Country = "FR"
	SP  Country = "SP"
)

func (c *Country) UnmarshalJSON(b []byte) error {
	var operation string
	if err := json.Unmarshal(b, &operation); err != nil {
		return errors.Wrap(err, "failed to unmarshal operation")
	}
	*c = Country(operation)

	return c.Validate()
}

func (c *Country) Validate() error {
	switch *c {
	case UK, LTU, FR, SP:
		return nil
	}

	return errors.Wrap(validator.New().Var(string(*c), "eq=UK|eq=LTU|eq=FR|eq=SP"), "validation error")
}

type User struct {
	FirstName string  `json:"first_name" binding:"min=1,max=75" validate:"required"`
	LastName  string  `json:"last_name" binding:"min=1,max=75" validate:"required"`
	Nickname  string  `json:"nickname" binding:"min=1,max=75" validate:"required"`
	Password  string  `json:"password" binding:"min=1,max=75" validate:"required"`
	Email     string  `json:"email" binding:"min=1,max=100" validate:"required"`
	Country   Country `json:"country"`
}
