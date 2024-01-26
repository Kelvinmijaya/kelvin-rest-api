package auth

import (
	"fmt"
	"net/http"
	"time"

	configs "github.com/Kelvinmijaya/kelvin-rest-api/config"
	"github.com/Kelvinmijaya/kelvin-rest-api/domain"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"
)

func GetJWTSecret() string {
	jwtSecretKey := configs.EnvConfigs.JWTSecret
	return jwtSecretKey
}

func GetRefreshJWTSecret() string {
	jwtRefreshSecretKey := configs.EnvConfigs.JWTRefreshToken
	return jwtRefreshSecretKey
}

// JWT Middleware Config to check valid token
func GetJWTMiddlewareConfig() echojwt.Config {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(Claims)
		},
		SigningKey:  []byte(GetJWTSecret()),
		TokenLookup: "cookie:" + accessTokenCookieName, // "<source>:<name>"
		ErrorHandler: func(c echo.Context, err error) error {
			cookie, errs := c.Cookie("access-token")
			if errs != nil {
				return errs
			}
			fmt.Println(cookie.Name)
			fmt.Println(cookie.Value)
			return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
		},
	}

	return config
}

// Function to clear out access token and refresh token
func LogoutTokenSetCookies(c echo.Context) (bool, error) {
	cookie := new(http.Cookie)
	cookie.Name = accessTokenCookieName
	cookie.Value = ""
	cookie.MaxAge = -1
	cookie.Path = "/"
	c.SetCookie(cookie)

	cookie.Name = refreshTokenCookieName
	c.SetCookie(cookie)

	return true, nil
}

// GenerateTokensAndSetCookies generates jwt token and saves it to the http-only cookie.
func GenerateTokensAndSetCookies(c echo.Context, us *domain.User) error {
	accessToken, exp, err := generateAccessToken(us)
	if err != nil {
		return err
	}

	setTokenCookie(accessTokenCookieName, accessToken, exp, c)

	// We generate here a new refresh token and saving it to the cookie.
	refreshToken, exp, err := generateRefreshToken(us)
	if err != nil {
		return err
	}

	//Refresh Token
	setTokenCookie(refreshTokenCookieName, refreshToken, exp, c)
	return nil
}

func generateAccessToken(us *domain.User) (string, time.Time, error) {
	// Declare the expiration time of the token (1h).
	expirationTime := time.Now().Add(1 * time.Hour)

	return generateToken(us, expirationTime, []byte(GetJWTSecret()))
}

func generateRefreshToken(us *domain.User) (string, time.Time, error) {
	// Declare the expiration time of the token - 24 hours.
	expirationTime := time.Now().Add(24 * time.Hour)

	return generateToken(us, expirationTime, []byte(GetRefreshJWTSecret()))
}

// Pay attention to this function. It holds the main JWT token generation logic.
func generateToken(us *domain.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	// Create the JWT claims, which includes the username and expiry time.
	claims := &Claims{
		Name: us.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Now(), err
	}

	return t, expirationTime, nil
}

// Here we are creating a new cookie, which will store the valid JWT token.
func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	// Http-only helps mitigate the risk of client side script accessing the protected cookie.
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

// JWTErrorChecker will be executed when user try to access a protected path.
func JWTErrorChecker(err error, c echo.Context) error {
	// Redirects to the signIn form.
	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("userSignInForm"))
}

// TokenRefresherMiddleware middleware, which refreshes JWT tokens if the access token is about to expire.
func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// If the user is not authenticated (no user token data in the context), don't do anything.
		if c.Get("user") == nil {
			c.Response().Writer.WriteHeader(http.StatusUnauthorized)
		}

		// Gets user token from the context.
		u := c.Get("user").(*jwt.Token)
		claims := u.Claims.(*Claims)
		expires := claims.RegisteredClaims.ExpiresAt

		// We ensure that a new token is not issued until enough time has elapsed.
		// In this case, a new token will only be issued if the old token is within
		// 15 mins of expiry.
		//lint:ignore S1024 ingore this for now
		if expires.Sub(time.Now()) < 15*time.Minute {
			// Gets the refresh token from the cookie.
			rc, err := c.Cookie(refreshTokenCookieName)
			if err == nil && rc != nil {
				// Parses token and checks if it valid.
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(GetRefreshJWTSecret()), nil
				})
				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						c.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}

				if tkn != nil && tkn.Valid {
					// If everything is good, update tokens.
					_ = GenerateTokensAndSetCookies(c, &domain.User{
						Name: claims.Name,
					})
				}
			}
		}

		return next(c)
	}
}
