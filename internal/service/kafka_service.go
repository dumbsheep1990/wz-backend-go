package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/go-redis/redis/v8"
)

// KafkaService 处理Kafka消息队列
type KafkaService interface {
	// 发布用户事件
	PublishUserEvent(ctx context.Context, eventType string, userID int64, data interface{}) error
	// 发布内容事件
	PublishContentEvent(ctx context.Context, eventType string, contentID int64, data interface{}) error
	// 消费用户事件
	ConsumeUserEvents(ctx context.Context, handler func(eventType string, userID int64, data json.RawMessage) error)
	// 消费内容事件
	ConsumeContentEvents(ctx context.Context, handler func(eventType string, contentID int64, data json.RawMessage) error)
}

// Event 事件结构
type Event struct {
	ID        string          `json:"id"`        // 事件唯一ID
	Type      string          `json:"type"`      // 事件类型
	EntityID  int64           `json:"entity_id"` // 实体ID（用户ID或内容ID）
	Data      json.RawMessage `json:"data"`      // 事件数据
	Timestamp time.Time       `json:"timestamp"` // 事件时间戳
}

type kafkaService struct {
	userEventProducer    *kafka.Writer
	contentEventProducer *kafka.Writer
	userEventConsumer    *kafka.Reader
	contentEventConsumer *kafka.Reader
	redis                *redis.Client
}

// NewKafkaService 创建Kafka服务
func NewKafkaService(
	userEventProducer *kafka.Writer,
	contentEventProducer *kafka.Writer,
	userEventConsumer *kafka.Reader,
	contentEventConsumer *kafka.Reader,
	redis *redis.Client,
) KafkaService {
	return &kafkaService{
		userEventProducer:    userEventProducer,
		contentEventProducer: contentEventProducer,
		userEventConsumer:    userEventConsumer,
		contentEventConsumer: contentEventConsumer,
		redis:                redis,
	}
}

// PublishUserEvent 发布用户事件
func (s *kafkaService) PublishUserEvent(ctx context.Context, eventType string, userID int64, data interface{}) error {
	return s.publishEvent(ctx, s.userEventProducer, eventType, userID, data)
}

// PublishContentEvent 发布内容事件
func (s *kafkaService) PublishContentEvent(ctx context.Context, eventType string, contentID int64, data interface{}) error {
	return s.publishEvent(ctx, s.contentEventProducer, eventType, contentID, data)
}

// 内部方法：发布事件
func (s *kafkaService) publishEvent(ctx context.Context, producer *kafka.Writer, eventType string, entityID int64, data interface{}) error {
	if producer == nil {
		return errors.New("producer is nil")
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	event := Event{
		ID:        generateUniqueID(),
		Type:      eventType,
		EntityID:  entityID,
		Data:      dataBytes,
		Timestamp: time.Now(),
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   []byte(eventType),
		Value: eventBytes,
	}

	// 事件去重检查 (使用Redis实现幂等性)
	dedupeKey := "event:" + event.ID
	_, err = s.redis.SetNX(ctx, dedupeKey, 1, 24*time.Hour).Result()
	if err != nil {
		return err
	}

	return producer.WriteMessages(ctx, message)
}

// ConsumeUserEvents 消费用户事件
func (s *kafkaService) ConsumeUserEvents(ctx context.Context, handler func(eventType string, userID int64, data json.RawMessage) error) {
	s.consumeEvents(ctx, s.userEventConsumer, handler)
}

// ConsumeContentEvents 消费内容事件
func (s *kafkaService) ConsumeContentEvents(ctx context.Context, handler func(eventType string, contentID int64, data json.RawMessage) error) {
	s.consumeEvents(ctx, s.contentEventConsumer, handler)
}

// 内部方法：消费事件
func (s *kafkaService) consumeEvents(ctx context.Context, consumer *kafka.Reader, handler func(eventType string, entityID int64, data json.RawMessage) error) {
	if consumer == nil {
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				message, err := consumer.ReadMessage(ctx)
				if err != nil {
					// 记录错误并继续
					continue
				}

				var event Event
				err = json.Unmarshal(message.Value, &event)
				if err != nil {
					// 记录错误并继续
					continue
				}

				// 事件去重检查 (使用Redis实现幂等性)
				dedupeKey := "event:" + event.ID
				exists, _ := s.redis.Exists(ctx, dedupeKey).Result()
				if exists > 0 {
					// 事件已处理，跳过
					continue
				}

				// 处理事件
				err = handler(event.Type, event.EntityID, event.Data)
				if err != nil {
					// 记录错误，可能需要重试或死信队列
					continue
				}

				// 标记事件为已处理
				_, _ = s.redis.SetEX(ctx, dedupeKey, 1, 24*time.Hour).Result()
			}
		}
	}()
}

// 生成唯一ID
func generateUniqueID() string {
	// 使用时间戳和随机数生成唯一ID
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond)
	}
	return string(b)
}
