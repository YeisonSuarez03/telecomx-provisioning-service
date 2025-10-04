package config

import (
	"os"
)

type Config struct {
	MongoURI string
	Brokers  []string
	Topic    string
	Group    string
	Client   string
	Port     string
}

var config *Config

func initEnvironment() {
	config = &Config{
		MongoURI: os.Getenv("MONGODB_URI"),
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    os.Getenv("KAFKA_TOPIC"),
		Group:    os.Getenv("KAFKA_GROUP_ID"),
		Client:   os.Getenv("KAFKA_CLIENT_ID"),
		Port:     os.Getenv("PORT"),
	}
}

func InstanceConfig() Config {
	if config == nil {
		initEnvironment()
	}
	return *config
}
