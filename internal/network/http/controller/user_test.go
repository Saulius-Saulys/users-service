package controller_test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/Saulius-Saulys/users-service/internal/messaging/rabbitmq"
	"github.com/Saulius-Saulys/users-service/internal/mocks"
	"github.com/Saulius-Saulys/users-service/internal/model"
	"github.com/Saulius-Saulys/users-service/internal/network/http/controller"
	"github.com/Saulius-Saulys/users-service/internal/network/http/controller/dto"
	"github.com/Saulius-Saulys/users-service/internal/repository"
	"github.com/Saulius-Saulys/users-service/internal/service"
	"github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var userController *controller.User
var rabbitMQPublisherMock *mocks.MockRabbitMQPublisher

func TestCreate(t *testing.T) {
	setup(t)
	router := userRouter(userController)

	t.Run("Create user success request", func(t *testing.T) {
		validUserDTO := createUserDTO("test@email.com", "LTU")
		r := gofight.New()
		r.POST("/user-service/users/").
			SetJSONInterface(validUserDTO).
			Run(router, func(response gofight.HTTPResponse, _ gofight.HTTPRequest) {
				require.Equal(t, http.StatusCreated, response.Code)
				userFromDB := model.User{}
				result := gormConn.Where("first_name = ? and email = ?", validUserDTO.FirstName, validUserDTO.Email).Find(&userFromDB)
				assert.Equal(t, int64(1), result.RowsAffected)
				assert.NoError(t, result.Error)

				// compare DTO from request and actual model from DB
				assertDTOToModel(t, validUserDTO, userFromDB)

				gormConn.Delete(userFromDB)
			})
	})

	t.Run("Create user country validation error", func(t *testing.T) {
		invalidUserCountryDTO := createUserDTO("test@email.com", "PL")
		r := gofight.New()
		r.POST("/user-service/users/").
			SetJSONInterface(invalidUserCountryDTO).
			Run(router, func(response gofight.HTTPResponse, _ gofight.HTTPRequest) {
				require.Equal(t, http.StatusUnprocessableEntity, response.Code)
			})
	})

	t.Run("Create user email validation error", func(t *testing.T) {
		invalidUserEmailDTO := createUserDTO("invalid email", string(dto.LTU))
		r := gofight.New()
		r.POST("/user-service/users/").
			SetJSONInterface(invalidUserEmailDTO).
			Run(router, func(response gofight.HTTPResponse, _ gofight.HTTPRequest) {
				require.Equal(t, http.StatusUnprocessableEntity, response.Code)
			})
	})
}

func TestUpdate(t *testing.T) {
	setup(t)
	userUUID, err := uuid.NewUUID()
	assert.NoError(t, err)
	originalModel := createUserModel(userUUID.String(), "test first name", dto.LTU)

	updatedFirstName := "updated name"
	updatedEmail := "updated@email.com"
	updateUserDTO := dto.UpdateUser{
		FirstName: &updatedFirstName,
		Email:     &updatedEmail,
	}

	t.Run("Update user", func(t *testing.T) {
		insertResult := gormConn.Create(&originalModel)
		assert.NoError(t, insertResult.Error)

		router := userRouter(userController)
		r := gofight.New()
		path := "/user-service/users/" + originalModel.ID
		rabbitMQPublisherMock.EXPECT().PublishMessage(mock.Anything, mock.Anything).Return(nil).Times(1)

		r.PUT(path).SetJSONInterface(updateUserDTO).
			Run(router, func(response gofight.HTTPResponse, _ gofight.HTTPRequest) {
				require.Equal(t, http.StatusOK, response.Code)
				updatedUserFromDB := model.User{}
				selectResult := gormConn.Where("first_name = ? and email = ?", updatedFirstName, updatedEmail).Find(&updatedUserFromDB)

				assert.Equal(t, int64(1), selectResult.RowsAffected)
				assert.NoError(t, selectResult.Error)

				assert.Equal(t, updatedFirstName, updatedUserFromDB.FirstName)
				assert.Equal(t, updatedEmail, updatedUserFromDB.Email)

				gormConn.Delete(originalModel)
			})
	})

	t.Run("Error from rabbitMQ received", func(t *testing.T) {
		router := userRouter(userController)
		r := gofight.New()
		path := "/user-service/users/" + originalModel.ID
		rabbitMQPublisherMock.EXPECT().PublishMessage(mock.Anything, mock.Anything).Return(errors.New("something went wrong")).Times(1)

		r.PUT(path).SetJSONInterface(updateUserDTO).
			Run(router, func(response gofight.HTTPResponse, _ gofight.HTTPRequest) {
				require.Equal(t, http.StatusBadRequest, response.Code)
			})
	})
}

func TestDelete(t *testing.T) {
	setup(t)
	router := userRouter(userController)

	userUUID, err := uuid.NewUUID()
	assert.NoError(t, err)
	originalModel := createUserModel(userUUID.String(), "test first name", dto.LTU)

	t.Run("Delete user", func(t *testing.T) {
		insertResult := gormConn.Create(&originalModel)
		assert.NoError(t, insertResult.Error)

		r := gofight.New()
		path := "/user-service/users/" + originalModel.ID
		rabbitMQPublisherMock.EXPECT().PublishMessage(
			mock.Anything,
			&rabbitmq.OutputMessage{
				Action: rabbitmq.ActionDelete,
				ID:     originalModel.ID,
			}).Return(nil).Times(1)

		r.DELETE(path).
			Run(router, func(response gofight.HTTPResponse, _ gofight.HTTPRequest) {
				require.Equal(t, http.StatusOK, response.Code)
				deleteUserFromDB := model.User{}
				selectResult := gormConn.Where("first_name = ? and email = ?", originalModel.FirstName, originalModel.Email).Find(&deleteUserFromDB)

				// user should be deleted
				assert.Equal(t, int64(0), selectResult.RowsAffected)
				assert.NoError(t, selectResult.Error)

				gormConn.Delete(originalModel)
			})
	})

	t.Run("Delete user", func(t *testing.T) {
		r := gofight.New()
		path := "/user-service/users/" + originalModel.ID
		rabbitMQPublisherMock.EXPECT().PublishMessage(
			mock.Anything,
			&rabbitmq.OutputMessage{
				Action: rabbitmq.ActionDelete,
				ID:     originalModel.ID,
			}).Return(errors.New("something went wrong")).Times(1)

		r.DELETE(path).
			Run(router, func(response gofight.HTTPResponse, _ gofight.HTTPRequest) {
				require.Equal(t, http.StatusBadRequest, response.Code)
			})
	})
}

func TestGetByCountry(t *testing.T) {
	setup(t)

	userUUID, err := uuid.NewUUID()
	assert.NoError(t, err)
	userUUID2, err := uuid.NewUUID()
	assert.NoError(t, err)
	userUUID3, err := uuid.NewUUID()
	assert.NoError(t, err)
	originalModel := createUserModel(userUUID.String(), "test first name", dto.LTU)
	originalModel2 := createUserModel(userUUID2.String(), "test first name 2", dto.UK)
	originalModel3 := createUserModel(userUUID3.String(), "test first name 3", dto.LTU)

	t.Run("Get user by country", func(t *testing.T) {
		insertResult := gormConn.Create(&originalModel)
		assert.NoError(t, insertResult.Error)

		insertResult2 := gormConn.Create(&originalModel2)
		assert.NoError(t, insertResult2.Error)

		insertResult3 := gormConn.Create(&originalModel3)
		assert.NoError(t, insertResult3.Error)

		router := userRouter(userController)
		r := gofight.New()
		path := "/user-service/users/" + "LTU"

		r.GET(path).
			Run(router, func(response gofight.HTTPResponse, _ gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, response.Code)
				body, readErr := io.ReadAll(response.Body)
				assert.NoError(t, readErr)
				ltuUsers := []model.User{originalModel, originalModel3}
				expectedResponse, marshalErr := json.Marshal(ltuUsers)
				assert.NoError(t, marshalErr)

				assert.Equal(t, expectedResponse, body)

				gormConn.Delete(originalModel)
				gormConn.Delete(originalModel2)
				gormConn.Delete(originalModel3)
			})
	})
}

func createUserDTO(email string, country string) dto.CreateUser {
	return dto.CreateUser{
		FirstName: "test name",
		LastName:  "test last name",
		Nickname:  "test nickname",
		Password:  "test_password",
		Email:     email,
		Country:   dto.Country(country),
	}
}

func createUserModel(id, firstName string, country dto.Country) model.User {
	return model.User{
		ID:        id,
		FirstName: firstName,
		LastName:  "test last name",
		Nickname:  "test nickname",
		Password:  "test_password",
		Email:     "test@email.com",
		Country:   country,
	}
}

func assertDTOToModel(t *testing.T, createUserDTO dto.CreateUser, userFromDB model.User) {
	assert.Equal(t, createUserDTO.FirstName, userFromDB.FirstName)
	assert.Equal(t, createUserDTO.LastName, userFromDB.LastName)
	assert.Equal(t, createUserDTO.Nickname, userFromDB.Nickname)
	assert.Equal(t, createUserDTO.Email, userFromDB.Email)
	assert.Equal(t, createUserDTO.Country, userFromDB.Country)
}

func userRouter(userController *controller.User) *gin.Engine {
	router := gin.New()
	userEndpoints := router.Group("/user-service/users/")
	userEndpoints.POST("/", userController.Create)
	userEndpoints.PUT("/:id", userController.Update)
	userEndpoints.DELETE("/:id", userController.Delete)
	userEndpoints.GET("/:country", userController.GetByCountry)

	return router
}

func setup(t *testing.T) {
	logger := zap.NewNop()
	userRepository := repository.NewUser(logger, gormConn)
	rabbitMQPublisherMock = mocks.NewMockRabbitMQPublisher(t)
	userService := service.NewUser(logger, userRepository, rabbitMQPublisherMock)
	userController = controller.NewUser(logger, userService)
}
