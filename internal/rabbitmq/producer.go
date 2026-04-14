package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQProducer(url string) (*RabbitMQProducer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare the exchange
	err = ch.ExchangeDeclare(
		"statistics-exchange", // name
		"topic",                // type
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// AUTOMATIC SETUP: Declare and bind the achievement-service-queue
	// This makes it work "automatically" as the user requested.
	_, err = ch.QueueDeclare(
		"achievement-service-queue", // name
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)
	if err != nil {
		log.Printf("Warning: failed to declare queue: %v", err)
	} else {
		err = ch.QueueBind(
			"achievement-service-queue", // queue name
			"user.action.#",             // routing key pattern
			"statistics-exchange",      // exchange
			false,
			nil,
		)
		if err != nil {
			log.Printf("Warning: failed to bind queue: %v", err)
		}
	}

	return &RabbitMQProducer{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *RabbitMQProducer) PublishEvent(routingKey string, event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(ctx,
		"statistics-exchange", // exchange
		routingKey,             // routing key
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		log.Printf("[RabbitMQ] Failed to publish event: %v", err)
		return err
	}

	return nil
}

func (p *RabbitMQProducer) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
