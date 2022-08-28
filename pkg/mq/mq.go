package mq

import (
	"context"

	"gitlab.com/g6834/team31/tasks/pkg/mq/kafka"
	"gitlab.com/g6834/team31/tasks/pkg/mq/types"
)

type Producer interface {
	SendMessage(ctx context.Context, messages []types.Message, partition int) error
}

func NewProducer(brokers []string, topic string) (Producer, error) {
	return kafka.NewProducer(brokers, topic)
}

type Consumer interface {
	GetMessage(ctx context.Context) (types.Message, error)
	ReadAndCommit(ctx context.Context, fn func(ctx context.Context, m types.Message) error) (types.Message, error)
}

func NewConsumer(brokers []string, topic string, groupId string) (Consumer, error) {
	return kafka.NewConsumer(brokers, topic, groupId)
}