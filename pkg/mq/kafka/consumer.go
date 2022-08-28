package kafka

import (
	"context"
	"errors"
	"fmt"

	kafka "github.com/segmentio/kafka-go"
	"gitlab.com/g6834/team31/tasks/pkg/mq/types"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

type Consumer struct {
	Reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) (*Consumer, error) {
	if len(brokers) == 0 || brokers[0] == "" || topic == "" {
		return nil, errors.New("нет парамемтров подкл к кафке")
	}

	c := &Consumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 1,
			MaxBytes: 10e6,
			// Partition: 0,
		}),
	}
	return c, nil
}

//GetMessage помечает прочитанное сообщение как обработанное и удаляет его из кафки
func (c *Consumer) GetMessage(ctx context.Context) (types.Message, error) {
	kafkaMsg, err := c.Reader.ReadMessage(ctx)
	if err != nil {
		return types.Message{}, fmt.Errorf("failed to read message: %w", err)
	}
	return types.Message{
		Key:   kafkaMsg.Key,
		Value: kafkaMsg.Value,
	}, err
}

func (c *Consumer) ReadAndCommit(ctx context.Context, fn func(ctx context.Context, m types.Message) error) (types.Message, error) {
	kafkaMsg, err := c.Reader.FetchMessage(ctx)
	ctx, span := otel.Tracer("team31").Start(ctx, "consumer readAndCommit kafka", trace.WithAttributes(semconv.MessagingOperationProcess))
	defer span.End()
	if err != nil {
		return types.Message{}, fmt.Errorf("failed to read message: %w", err)
	}
	m := types.Message{
		Key:   kafkaMsg.Key,
		Value: kafkaMsg.Value,
	}
	if fn != nil {
		if err := fn(ctx, m); err != nil {
			return types.Message{}, fmt.Errorf("failed to process message: %w", err)
		}
	}
	if err := c.Reader.CommitMessages(ctx, kafkaMsg); err != nil {
		return types.Message{}, fmt.Errorf("failed to commit the read message: %w", err)
	}
	return m, nil
}
