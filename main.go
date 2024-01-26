package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	configs "github.com/Kelvinmijaya/kelvin-rest-api/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"

	_articleHttpDelivery "github.com/Kelvinmijaya/kelvin-rest-api/article/delivery/http"
	_articleRepo "github.com/Kelvinmijaya/kelvin-rest-api/article/repository/postgres"
	_articleUsecase "github.com/Kelvinmijaya/kelvin-rest-api/article/usecase"

	_userHttpDelivery "github.com/Kelvinmijaya/kelvin-rest-api/user/delivery/http"
	_userRepo "github.com/Kelvinmijaya/kelvin-rest-api/user/repository/postgres"
	_userUsecase "github.com/Kelvinmijaya/kelvin-rest-api/user/usecase"
)

func init() {
	// Config ENV
	configs.InitEnvConfigs()
}

func main() {
	// DB Connection
	// Construct the full connection string
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		configs.EnvConfigs.DBHost, configs.EnvConfigs.DBPort, configs.EnvConfigs.DBUser, configs.EnvConfigs.DBPassword, configs.EnvConfigs.DBName, configs.EnvConfigs.DBssl)

	// Connect to the Cloud SQL database
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Echo Framework
	e := echo.New()
	// Init Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"}, // Replace with your frontend URL
		AllowCredentials: true,
	}))
	timeoutContext := time.Duration(configs.EnvConfigs.Timeout) * time.Second

	// Init Default
	e.GET("/", func(c echo.Context) error {
		cookie, err := c.Cookie("access-token")
		if err != nil {
			return err
		}
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
		return c.HTML(http.StatusOK, "Hello, World!")
	})

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	// Init Article
	ar := _articleRepo.NewPostgresArticleRepository(db)
	au := _articleUsecase.NewArticleUsecase(ar, timeoutContext)
	_articleHttpDelivery.NewArticleHandler(e, au)

	// Init User
	ur := _userRepo.NewPostgresUserRepository(db)
	uu := _userUsecase.NewUserUsecase(ur, timeoutContext)
	_userHttpDelivery.NewUserHandler(e, uu)

	// Setup Server Address
	e.Logger.Fatal(e.Start(":" + string(configs.EnvConfigs.Port)))
}
