package sdkrabbitmq

import (
	"context"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

/**
 * rabbitmq_operations
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

const DefaultTimeout = 5 * time.Second

type RabbitData struct {
	Exchange     string
	ExchangeType string
	RoutingKey   string
	Payload      []byte
	Ctx          context.Context
}

func (p *Rabbitmq) Publisher(data RabbitData) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	ctx, cancel := getContext(data.Ctx)
	defer cancel()

	if err := p.channel.ExchangeDeclare(data.Exchange,
		data.ExchangeType,
		true,
		false,
		false,
		false,
		nil); err != nil {
		message = fmt.Sprintf("Failed to declare direct exchange (%s)", data.Exchange)
		log.Println(message)
		return err
	}

	err := p.channel.PublishWithContext(
		ctx,
		data.Exchange,
		data.RoutingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			DeliveryMode:    amqp.Transient,
			Body:            data.Payload,
			Priority:        0,
		},
	)
	if err != nil {
		message = fmt.Sprintf("Failed to publish a message: %s", err)
		return errors.New(message)
	}

	return nil
}

func (p *Rabbitmq) Subscriber(data RabbitData) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	//ctx, cancel := getContext(data.Ctx)
	//defer cancel()

	if err := p.channel.ExchangeDeclare(
		data.Exchange,
		data.ExchangeType,
		false,
		false,
		false,
		false,
		nil); err != nil {
		message = fmt.Sprintf("Failed to declare direct exchange (%s)", data.Exchange)
		log.Println(message)
		return
	}

	_, err := p.channel.QueueDeclare(
		p.queue,
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		message = fmt.Sprintf("Failed to declare a queue: %s", err)
		log.Println(message)
		return
	}

	if err := p.channel.QueueBind(
		p.queue,
		data.RoutingKey,
		data.Exchange,
		false,
		nil,
	); err != nil {
		message = fmt.Sprintf("Failed to bind a queue: %s", err)
		log.Fatal(message)
	}

	msg, err := p.channel.Consume(p.queue,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		message = fmt.Sprintf("Failed to register a consumer: %s", err)
		log.Fatal(message)
	}

	for d := range msg {
		p.handler(d.Body)
		err := d.Ack(false)
		if err != nil {
			message = fmt.Sprintf("Failed to ack: %s", err)
			log.Fatal(message)
			return
		}
	}
}

func getContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		return context.WithTimeout(context.Background(), DefaultTimeout)
	}
	return context.WithTimeout(ctx, DefaultTimeout)
}
