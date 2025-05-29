package eventbus

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"

	"wz-backend-go/internal/domain/user/repository"
)

// RabbitMQEventPublisher RabbitMQ事件发布器
type RabbitMQEventPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewRabbitMQEventPublisher 创建RabbitMQ事件发布器
func NewRabbitMQEventPublisher(amqpURL string) (repository.EventPublisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		"domain_events", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &RabbitMQEventPublisher{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

// Publish 发布事件
func (p *RabbitMQEventPublisher) Publish(event interface{}) error {
	// 序列化事件
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 发布事件
	err = p.channel.Publish(
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// Close 关闭连接
func (p *RabbitMQEventPublisher) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			log.Printf("failed to close channel: %v", err)
		}
	}
	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}
	return nil
}

// LogEventPublisher 日志事件发布器
type LogEventPublisher struct{}

// NewLogEventPublisher 创建日志事件发布器
func NewLogEventPublisher() repository.EventPublisher {
	return &LogEventPublisher{}
}

// Publish 发布事件
func (p *LogEventPublisher) Publish(event interface{}) error {
	// 序列化事件
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// 记录事件
	log.Printf("Event published: %s", string(body))

	return nil
}
