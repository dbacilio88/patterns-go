package app

import (
	"errors"
	"fmt"
	config "github.com/dbacilio88/patterns-go/internal/config/app"
)

/**
 *
 * THIS COMPONENT WAS BUILT FOLLOWING THE DEVELOPMENT STANDARDS AND PROCEDURE
 * NOVOPAYMENT APP DEVELOPMENT AND IS PROTECTED BY THE LAWS OF
 * INTELLECTUAL PROPERTY AND COPYRIGHT.
 *
 * @author cbacilio
 * @since 15/07/2025
 */

var message string

func ConfigureApplication(configPath string) error {
	var err error
	if err = loadRabbitConfig(); err != nil {
		message = fmt.Sprintf("Load Rabbit Config Error: %s", err.Error())
		return errors.New(message)
	}
	return nil
}

func GetConfig() *config.Config {
	if config.GlobalConfig == nil {
		message = fmt.Sprintf("Load error config")
		return nil
	}
	return config.GlobalConfig
}

func GetConfigGlobal(key string) string {
	return config.Secret[key]
}

func loadRabbitConfig() error {
	cfg := GetConfig()
	if cfg == nil {
		message = fmt.Sprintf("config must be initialized before rabbitmq configurer")
		return errors.New(message)
	}

	cfg.RabbitMQ.URI = "amqp://guest:guest@localhost:5672/"

	return nil
}

func ExecuteRabbitProcess() error {

	return nil
}
