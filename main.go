package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	_articleHttpDelivery "github.com/Kelvinmijaya/kelvin-rest-api/article/delivery/http"
	_articleHttpDeliveryMiddleware "github.com/Kelvinmijaya/kelvin-rest-api/article/delivery/http/middleware"
	_articleRepo "github.com/Kelvinmijaya/kelvin-rest-api/article/repository/postgres"
	_articleUsecase "github.com/Kelvinmijaya/kelvin-rest-api/article/usecase"

	_userHttpDelivery "github.com/Kelvinmijaya/kelvin-rest-api/user/delivery/http"
	_userRepo "github.com/Kelvinmijaya/kelvin-rest-api/user/repository/postgres"
	_userUsecase "github.com/Kelvinmijaya/kelvin-rest-api/user/usecase"
)

func init() {
	// Set viper path and read configuration
	viper.SetConfigType("json")
	viper.AddConfigPath("./conf")
	if os.Getenv("ENV") == "PRODUCTION" {
		viper.SetConfigName("production")
	} else {
		viper.SetConfigName("development")
	}
	err := viper.ReadInConfig()

	// Handle errors reading the config file
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		log.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	// DB Init
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	dbSSLmode := viper.GetString(`database.sslmode`)
	fmt.Sprintln(dbName)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, dbSSLmode)

	db, err := sql.Open("postgres", psqlInfo)

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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	middL := _articleHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	// Init Default
	e.GET("/", func(c echo.Context) error {
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
	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "9090"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
