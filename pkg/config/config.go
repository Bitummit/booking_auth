package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)


type (
	Config struct {
		Env string `yaml:"env" env-default:"dev"`
		HTTPServer `yaml:"http_server"`
		GrpcServer `yaml:"grpc_server"`
		KafkaServer `yaml:"kafka"`
	}

	HTTPServer struct {
		HTTPAddress string `yaml:"http_address" env-default:"0.0.0.0:8000"`
	}

	GrpcServer struct {
		GrpcAddress string `yaml:"grpc_address"`
	}

	KafkaServer struct {
		KafkaBrokerAddress string `yaml:"kafka_broker_address"`
	}
)

func InitConfig() *Config{
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file!")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalln("Empty path")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalln("Can not find config file")
	} 
	
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalln("Error in reading config file!")
	}
	
	return &cfg
}