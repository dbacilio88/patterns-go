package sdkrabbitmq

import (
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

/**
 * rabbitmq_connection
 * <p>
 * This file contains core data structures and logic used throughout the application.
 *
 * <p><strong>Copyright © 2025 – All rights reserved.</strong></p>
 *
 * <p>This source code is distributed under a collaborative license.</p>
 *
 * <p>
 * Contributions, suggestions, and improvements are welcome!
 * You are free to fork, modify, and submit pull requests under the terms of the repository's license.
 * Please ensure proper attribution to the original author(s) and preserve this notice in derivative works.
 * </p>
 *
 * @author Christian Bacilio De La Cruz
 * @email dbacilio88@outlook.es
 * @since 7/11/2025
 */

type Rabbitmq struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	mutex      *sync.Mutex
	queue      string
	handler    func([]byte)
	IsExchange bool
}

type RabbitParam struct {
	Url       string //"amqp://guest:guest@localhost:5672/"
	QueueName string
	Vhost     string
	Handler   func([]byte)
}

var instance *Rabbitmq
var once sync.Once
var message string

func GetInstance() *Rabbitmq {
	once.Do(func() {
		instance = &Rabbitmq{}
	})
	return instance
}

func (p *Rabbitmq) Connect(param RabbitParam) error {
	var err error
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.conn, err = amqp.Dial(param.Url)

	if err != nil {
		message = fmt.Sprintf("Failed to connect to RabbitMQ. Url=%s Error=%s\n", param.Url, err.Error())
		return errors.New(message)
	}

	p.channel, err = p.conn.Channel()
	if err != nil {
		message = fmt.Sprintf("Failed to open a channel. Url=%s Error=%s\n", param.Url, err.Error())
		return errors.New(message)
	}

	table := amqp.Table{
		"x-queue-name":           param.QueueName,
		"x-deduplication-window": 60000,
	}

	_, err = p.channel.QueueDeclare(
		param.QueueName,
		true,
		false,
		false,
		false,
		table)

	if err != nil {
		message = fmt.Sprintf("Failed to declare a queue. Url=%s Error=%s\n", param.Url, err.Error())
		return errors.New(message)
	}

	p.queue = param.QueueName
	p.handler = param.Handler

	return nil
}

func (p *Rabbitmq) handleReconnection(param RabbitParam) {
	chanErr := p.conn.NotifyClose(make(chan *amqp.Error))
	if chanErr != nil {
		for {
			err := p.Connect(param)
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func (p *Rabbitmq) ConnectionWithHandler(param RabbitParam) error {
	var err error
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.conn, err = amqp.Dial(param.Url)
	if err != nil {
		message = fmt.Sprintf("Failed to open a channel. Url=%s Error=%s\n", param.Url, err.Error())
		return errors.New(message)
	}
	p.channel, err = p.conn.Channel()
	if err != nil {
		message = fmt.Sprintf("Failed to open a channel. Url=%s Error=%s\n", param.Url, err.Error())
		return errors.New(message)
	}
	p.handler = param.Handler
	go p.handlerReconnectionWithHandler(param)
	return nil
}

func (p *Rabbitmq) handlerReconnectionWithHandler(param RabbitParam) {
	chanErr := p.conn.NotifyClose(make(chan *amqp.Error))
	if chanErr != nil {
		for {
			err := p.ConnectionWithHandler(param)
			if err == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}
	}
}
