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
	DBUrl           string `mapstructure:"DB_URL"`
}

// Initilize this variable to access the env values
var EnvConfigs *envConfigs

// We will call this in main.go to load the env variables
func InitEnvConfigs() {
	EnvConfigs = loadEnvVariables()
}

// Call to load the variables from env
func loadEnvVariables() (config *envConfigs) {
	env := "DEVELOPMENT"
	if envOS := os.Getenv("ENV"); envOS != "" {
		env = envOS
	}

	if env == "PRODUCTION" {
		viper.AutomaticEnv() // Read environment variables automatically

	} else {
		// Tell viper the path/location of your env file. If it is root just add "."
		viper.AddConfigPath(".")

		// Tell viper the name of your file
		viper.SetConfigFile(".env")

		// Tell viper the type of your file
		viper.SetConfigType("env")

		// Viper reads all the variables from env file and log error if any found
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("Error reading env file", err)
		}
	}

	// Viper unmarshals the loaded env varialbes into the struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	return
}
