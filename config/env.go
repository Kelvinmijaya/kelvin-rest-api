package configs

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// struct to map env values
type envConfigs struct {
	Environtment    string `mapstructure:"ENV"`
	Port            string `mapstructure:"PORT"`
	Timeout         int    `mapstructure:"TIMEOUT"`
	JWTSecret       string `mapstructure:"JWT_SECRET"`
	JWTRefreshToken string `mapstructure:"JWT_REFRESH_SECRET"`
	DBHost          string `mapstructure:"DB_HOST"`
	DBPort          string `mapstructure:"DB_PORT"`
	DBUser          string `mapstructure:"DB_USER"`
	DBPassword      string `mapstructure:"DB_PASS"`
	DBName          string `mapstructure:"DB_NAME"`
	DBssl           string `mapstructure:"DB_SSL"`
}

// Initilize this variable to access the env values
var EnvConfigs *envConfigs

// We will call this in main.go to load the env variables
func InitEnvConfigs() {
	EnvConfigs = loadEnvVariables()
}

// Call to load the variables from env
func loadEnvVariables() (config *envConfigs) {
	v := viper.New()
	env := "DEVELOPMENT"
	if envOS := os.Getenv("ENV"); envOS != "" {
		env = envOS
	}

	if env == "PRODUCTION" {
		v.AutomaticEnv() // Read environment variables automatically

		v.BindEnv("ENV")
		v.BindEnv("PORT")
		v.BindEnv("TIMEOUT")
		v.BindEnv("JWT_SECRET")
		v.BindEnv("JWT_REFRESH_SECRET")
		v.BindEnv("DB_HOST")
		v.BindEnv("DB_PORT")
		v.BindEnv("DB_USER")
		v.BindEnv("DB_PASS")
		v.BindEnv("DB_NAME")
		v.BindEnv("DB_SSL")

	} else {
		// Tell viper the path/location of your env file. If it is root just add "."
		v.AddConfigPath(".")

		// Tell viper the name of your file
		v.SetConfigFile(".env")

		// Tell viper the type of your file
		v.SetConfigType("env")

		// Viper reads all the variables from env file and log error if any found
		if err := v.ReadInConfig(); err != nil {
			log.Fatal("Error reading env file", err)
		}
	}

	// Viper unmarshals the loaded env varialbes into the struct
	if err := v.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	return
}
