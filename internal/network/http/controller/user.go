package controller

import (
	"net/http"
	"strconv"

	"github.com/Saulius-Saulys/users-service/internal/model"
	"github.com/Saulius-Saulys/users-service/internal/network/http/controller/dto"
	"github.com/Saulius-Saulys/users-service/internal/network/http/validation"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	idParam      = "id"
	countryParam = "country"
	pageQuery    = "page"
	defaultPage  = "1"
	limitQuery   = "limit"
	defaultLimit = "10"
)

type UserService interface {
	Create(user *dto.CreateUser) (*model.User, error)
	Update(id string, user *dto.UpdateUser) (*model.User, error)
	Delete(id string) error
	GetByCountry(country dto.Country, page, limit int) ([]model.User, error)
}

type User struct {
	logger      *zap.Logger
	userService UserService
}

func NewUser(logger *zap.Logger, userService UserService) *User {
	return &User{
		logger:      logger,
		userService: userService,
	}
}

// Create updates user by provided data
//
// @Schemes
// @Summary		Creates user.
// @Description	Creates user with provided from customer.
// @Tags		User
// @Param request body dto.CreateUser true "query params"
// @Success		200			{object} model.User	"Successfully user created."
// @Failure 	400			"Error appeared while creating user model."
// @Failure 	422			"Validation error."
// @Failure 	500			"Internal server error."
// @Router		/user-service/users [post]
func (u *User) Create(ctx *gin.Context) {
	reqDTO := &dto.CreateUser{}
	if errStatusCode, err := validation.ValidateJSONBody(ctx, reqDTO); err != nil {
		abortErr := ctx.AbortWithError(errStatusCode, err)
		if abortErr != nil {
			u.logger.Error("failed to send error response")
		}
		return
	}
	user, err := u.userService.Create(reqDTO)
	if err != nil {
		abortErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if abortErr != nil {
			u.logger.Error("failed to send error response")
		}
	}
	ctx.JSON(http.StatusCreated, user)
}

// Update updates user by provided data
//
// @Schemes
// @Summary		Updates user.
// @Description	Updates the user using provided data from customer.
// @Tags		User
// @Param		id	path	string 	true	"ID of user."
// @Param request body dto.UpdateUser true "query params"
// @Success		200			{object} model.User	"Successfully retrieved information of user."
// @Failure 	400			"Error appeared while updating user model."
// @Failure 	422			"Validation error."
// @Failure 	500			"Internal server error."
// @Router		/user-service/users/{id} [put]
func (u *User) Update(ctx *gin.Context) {
	userID := ctx.Param(idParam)
	reqDTO := &dto.UpdateUser{}
	if errStatusCode, err := validation.ValidateJSONBody(ctx, reqDTO); err != nil {
		abortErr := ctx.AbortWithError(errStatusCode, err)
		if abortErr != nil {
			u.logger.Error("failed to send error response")
		}
		return
	}

	user, err := u.userService.Update(userID, reqDTO)
	if err != nil {
		abortErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if abortErr != nil {
			u.logger.Error("failed to send error response")
		}
	}

	ctx.JSON(http.StatusOK, user)
}

// Delete deletes user
//
// @Schemes
// @Summary		deletes user.
// @Description	Deletes the user with ID provided.
// @Tags		User
// @Param		id	path	string 	true	"ID of user."
// @Success		200			"Successfully retrieved information of user."
// @Failure 	400			"Error appeared while deleting user model."
// @Failure 	500			"Internal server error."
// @Router		/user-service/users/{id} [delete]
func (u *User) Delete(ctx *gin.Context) {
	userID := ctx.Param(idParam)
	err := u.userService.Delete(userID)
	if err != nil {
		abortErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if abortErr != nil {
			u.logger.Error("failed to send error response")
		}
	}

	ctx.JSON(http.StatusOK, nil)
}

// GetByCountry get users by country
//
// @Schemes
// @Summary		Gets users by country.
// @Description	Gets users by country provided by customer.
// @Tags		User
// @Param		country	path	string 	true	"Country of user."
// @Param		page	query	string 	true	"Page to retrieve."
// @Param		limit	query	string 	true	"Limit of entries to retrieve."
// @Success		200			{object} []model.User	"Successfully retrieved information of user from specific country."
// @Failure 	400			"Error appeared while retrieving users."
// @Failure 	500			"Internal server error."
// @Router		/user-service/users/{country} [get]
func (u *User) GetByCountry(ctx *gin.Context) {
	params, err := u.retrieveGetByCountryParams(ctx)
	if err != nil {
		abortErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if abortErr != nil {
			u.logger.Error("failed to send error response")
		}
		return
	}
	users, err := u.userService.GetByCountry(dto.Country(params.country), params.page, params.limit)
	if err != nil {
		abortErr := ctx.AbortWithError(http.StatusBadRequest, err)
		if abortErr != nil {
			u.logger.Error("failed to send error response")
		}
	}

	ctx.JSON(http.StatusOK, users)
}

type getByCountryParams struct {
	country string
	page    int
	limit   int
}

func (u *User) retrieveGetByCountryParams(ctx *gin.Context) (*getByCountryParams, error) {
	country := ctx.Param(countryParam)
	page := ctx.DefaultQuery(pageQuery, defaultPage)
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert page to integer")
	}
	limit := ctx.DefaultQuery(limitQuery, defaultLimit)
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert limit to integer")
	}

	return &getByCountryParams{
		country: country,
		page:    pageInt,
		limit:   limitInt,
	}, nil
}
