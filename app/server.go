package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	_articleHttpDelivery "github.com/Kelvinmijaya/kelvin-rest-api/article/delivery/http"
	_articleHttpDeliveryMiddleware "github.com/Kelvinmijaya/kelvin-rest-api/article/delivery/http/middleware"
	_articleRepo "github.com/Kelvinmijaya/kelvin-rest-api/article/repository/postgres"
	_articleUcase "github.com/Kelvinmijaya/kelvin-rest-api/article/usecase"
)

func init() {
	// Config Init
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
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
	fmt.Sprintln(dbName)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
    dbHost, dbPort, dbUser, dbPass, dbName)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("DB Successfully connected!")

	// Echo Framework
	e := echo.New()
	// Init Middleware
	middL := _articleHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	
	// Init Article
	ar := _articleRepo.NewPostgresArticleRepository(db)
	au := _articleUcase.NewArticleUsecase(ar, timeoutContext)
	_articleHttpDelivery.NewArticleHandler(e, au)

	// Init User

	//Init Default
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	
	// Setup Server Address
	serverAddr := viper.GetString(`server.address`)
	e.Logger.Fatal(e.Start(serverAddr))
}