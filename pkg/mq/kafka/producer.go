package kafka

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/Shopify/sarama"
	kafka "github.com/segmentio/kafka-go"
	"gitlab.com/g6834/team31/tasks/pkg/mq/types"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	if len(brokers) == 0 || brokers[0] == "" || topic == "" {
		return nil, errors.New("нет парамемтров подкл к кафке")
	}

	p := &Producer{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers[0]),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
	return p, nil
}

func (p *Producer) SendMessage(ctx context.Context, messages []types.Message, partition int) error {
	for _, m := range messages {
		ctx, span := otel.Tracer("team31_tasks").Start(ctx, "service.CreateTask")
		defer span.End()
		// otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
		// var hdrs []kafka.Header
		// c := newKafkaHeaderCarrier(hdrs)
		msg := kafka.Message{
			Key:       m.Key,
			Value:     m.Value,
			Partition: partition,
			// Headers:   a, // добавляем спан
		}
		if err := p.Writer.WriteMessages(ctx, msg); err != nil {
			return fmt.Errorf("failer to send message to Kafka %w", err)
		}
		msg2 := sarama.ProducerMessage{
			Topic: "mytopic",
			Key:   sarama.StringEncoder("random_number"),
			Value: sarama.StringEncoder(fmt.Sprintf("%d", rand.Intn(1000))),
		}
		otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(&msg2))
		// otel.GetTextMapPropagator().Extract(ctx, otelsarama.NewConsumerMessageCarrier())
	}
	return nil
}
