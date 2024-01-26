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

	r := e.Group("/auth")
	// Restricted group
	r.Use(echojwt.WithConfig(auth.GetJWTMiddlewareConfig()))
	r.Use(auth.TokenRefresherMiddleware)
	r.GET("", handler.Auth)

	e.POST("/login", handler.Login)
	e.POST("/logout", handler.Logout)
}

// Login will try to check the user is valid or not
func (a *UserHandler) Login(c echo.Context) (err error) {
	var user domain.User
	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	ctx := c.Request().Context()

	//TODO for register logic
	// generate Hash
	// bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	// fmt.Println(string(bytes))

	// Usecase Login
	err = a.UUsecase.Login(ctx, &user)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	// Generate Token & Set Cookies
	err = auth.GenerateTokensAndSetCookies(c, &user)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, ResponseSuccess{Message: "Success Sign in!"})
}

// Login will try to check the user is valid or not
func (a *UserHandler) Auth(c echo.Context) (err error) {

	return c.JSON(http.StatusOK, ResponseSuccess{Message: "Authorized"})
}

func (a *UserHandler) Logout(c echo.Context) (err error) {
	_, err = auth.LogoutTokenSetCookies(c)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, ResponseSuccess{Message: "Success Logout!"})
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
