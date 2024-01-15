package http

import (
	"net/http"

	"github.com/Kelvinmijaya/kelvin-rest-api/domain"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ResponseSuccess struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UUsecase domain.UserUsecase
}

func NewUserHandler(e *echo.Echo, uu domain.UserUsecase) {
	handler := &UserHandler{
		UUsecase: uu,
	}

	e.POST("/login", handler.Login)
}

// Login will try to check the user is valid or not
func (a *UserHandler) Login(c echo.Context) (err error) {
	email := c.Param("email")
	password := c.Param("password")

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	// TODO: check valid email and password
	// var ok bool
	// if ok, err = isRequestValid(&user); !ok {
	// 	return c.JSON(http.StatusBadRequest, err.Error())
	// }

	ctx := c.Request().Context()
	err = a.UUsecase.Login(ctx, email, password)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, ResponseSuccess{Message: "success sign in"})
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case domain.ErrInternalServerError:
		return http.StatusInternalServerError
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// func isRequestValid(m *domain.User) (bool, error) {
// 	validate := validator.New()
// 	err := validate.Struct(m)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }
