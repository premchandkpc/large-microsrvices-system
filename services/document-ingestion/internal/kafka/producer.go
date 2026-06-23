package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
	"go.uber.org/zap"
)

type Producer struct {
	writer *kafka.Writer
	logger *zap.Logger
}

func NewProducer(brokers []string, logger *zap.Logger) (*Producer, error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.RoundRobin{},
		Compression:  compress.Snappy,
		BatchSize:    100,
		BatchTimeout: 10,
		RequiredAcks: kafka.RequireAll,
		Logger:       kafka.LoggerFunc(logger.Sugar().Infof),
		ErrorLogger:  kafka.LoggerFunc(logger.Sugar().Errorf),
	}

	return &Producer{writer: writer, logger: logger}, nil
}

func (p *Producer) Publish(ctx context.Context, topic string, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshaling message: %w", err)
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
		Headers: []kafka.Header{
			{Key: "content-type", Value: []byte("application/json")},
			{Key: "source", Value: []byte("document-ingestion")},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("writing message to kafka: %w", err)
	}

	p.logger.Info("message published",
		zap.String("topic", topic),
		zap.String("key", key),
	)
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
