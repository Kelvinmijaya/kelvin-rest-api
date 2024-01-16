package http

import (
	"net/http"

	"github.com/Kelvinmijaya/kelvin-rest-api/domain"
	auth "github.com/Kelvinmijaya/kelvin-rest-api/user/delivery/http/middleware"
	echojwt "github.com/labstack/echo-jwt/v4"
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

	// Restricted group
	r := e.Group("/r")
	r.Use(echojwt.WithConfig(auth.GetJWTMiddlewareConfig()))
	// r.Use(auth.TokenRefresherMiddleware)
	r.GET("", restricted)
}

// Test only
func restricted(c echo.Context) error {
	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(*jwtCustomClaims)
	// name := claims.Name
	return c.String(http.StatusOK, "Welcome !")
}

// Login will try to check the user is valid or not
func (a *UserHandler) Login(c echo.Context) (err error) {
	var user domain.User
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	// TODO: check valid email and password
	// var ok bool
	// if ok, err = isRequestValid(&user); !ok {
	// 	return c.JSON(http.StatusBadRequest, err.Error())
	// }

	ctx := c.Request().Context()
	err = a.UUsecase.Login(ctx, email, password, &user)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	// Generate Token & Set Cookies
	err = auth.GenerateTokensAndSetCookies(c, &user)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, ResponseSuccess{Message: "Success Sign in!"})
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
