package boot

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	validator "gopkg.in/go-playground/validator.v9"
)

var config Config

// Config contains most of the application's configuration objects
type Config struct {
	ServiceHost         string `json:"serviceHost" bson:"serviceHost" validate:"required"`
	MongoDatabaseName   string `json:"mongoDatabaseName" bson:"mongoDatabaseName" validate:"required"`
	MongoCollectionName string `json:"mongoCollectionName" bson:"mongoCollectionName" validate:"required"`
	MongoHost           string `json:"mongoHost" bson:"mongoHost" validate:"required"`
	MongoPort           string `json:"mongoPort" bson:"mongoPort" validate:"required"`
	MongoUserName       string `json:"mongoUserName" bson:"mongoUserName" validate:"required"`
	MongoPassword       string `json:"mongoPassword" bson:"mongoPassword" validate:"required"`
	MongoTimeout        int64  `json:"mongoTimeout" bson:"mongoTimeout" validate:"required"`
	ContextTimeout      int64  `json:"contextTimeout" bson:"contextTimeout" validate:"required"`
	JwtSecret           string `json:"jwtSecret" bson:"jwtSecret" validate:"required"`
}

// IsConfigValid checks for if the config is valid
func IsConfigValid(c *Config) (bool, error) {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetConfig returns the configuration of the application
func GetConfig() *Config {
	return &config
}

// SetupConfig sets up the config and send the config back to the application
func SetupConfig() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config.ServiceHost = os.Getenv("ServiceHost")
	config.MongoCollectionName = os.Getenv("MongoCollectionName")
	config.MongoDatabaseName = os.Getenv("MongoDatabaseName")
	config.MongoHost = os.Getenv("MongoHost")
	config.MongoPassword = os.Getenv("MongoPassword")
	config.MongoPort = os.Getenv("MongoPort")
	config.MongoUserName = os.Getenv("MongoUserName")
	MongoTimeout, err := strconv.ParseInt(os.Getenv("MongoTimeout"), 10, 64)
	config.MongoTimeout = MongoTimeout
	ContextTimeout, err := strconv.ParseInt(os.Getenv("ContextTimeout"), 10, 64)
	config.ContextTimeout = ContextTimeout
	config.JwtSecret = os.Getenv("JwtSecret")

	if ok, err := IsConfigValid(&config); !ok {
		return nil, err
	}
	return &config, nil
}
