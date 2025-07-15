package app

import "sync"

/**
 *
 * THIS COMPONENT WAS BUILT FOLLOWING THE DEVELOPMENT STANDARDS AND PROCEDURE
 * NOVOPAYMENT APP DEVELOPMENT AND IS PROTECTED BY THE LAWS OF
 * INTELLECTUAL PROPERTY AND COPYRIGHT.
 *
 * @author cbacilio
 * @since 15/07/2025
 */

type Config struct {
	RabbitMQ RabbitMQ `mapstructure:"rabbitmq"`
}

type RabbitMQ struct {
	URI string `mapstructure:"uri"`
}

var (
	GlobalConfig *Config
	Once         sync.Once
	Secret       map[string]string
)
