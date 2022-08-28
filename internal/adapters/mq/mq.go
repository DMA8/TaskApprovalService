package mq

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"gitlab.com/g6834/team31/tasks/pkg/grpc_task"
	"gitlab.com/g6834/team31/auth/pkg/logging"
	"gitlab.com/g6834/team31/tasks/pkg/mq/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	"gitlab.com/g6834/team31/tasks/internal/domain/models"
	"gitlab.com/g6834/team31/tasks/pkg/mq"
)

const (
	partition = 0 //??
)

type MQClient struct {
	ConsumerMailTopic mq.Consumer
	ProducerMailTopic mq.Producer
	ProducerTaskTopic mq.Producer
	l                 *logging.Logger
}

func NewMQClient(cMail mq.Consumer, pMail, pTasks mq.Producer, l *logging.Logger) *MQClient {
	return &MQClient{
		ConsumerMailTopic: cMail,
		ProducerMailTopic: pMail,
		ProducerTaskTopic: pTasks,
		l:                 l,
	}
}

func (m *MQClient) PushMail(ctx context.Context, mail *models.Mail) (models.TaskResponse, error) {
	value := grpc_task.Mail{
		Header:    mail.Body,
		Body:      mail.Body,
		CreateTs:  timestamppb.New(mail.CreateTS),
		EmailList: mail.EmailList,
	}
	valueBytes, err := proto.Marshal(&value)
	if err != nil {
		m.l.Warn().Err(err).Msg("MQClient.PushMail proto marshal problem")
		return models.TaskResponse{}, fmt.Errorf("couldn't encode mail in proto: %w", err)
	}
	msg := types.Message{
		Key:   []byte("push_mail tasks"),
		Value: valueBytes,
	}
	if err := m.ProducerMailTopic.SendMessage(ctx, []types.Message{msg}, partition); err != nil {
		m.l.Warn().Err(err).Msg("MQClient.PushMail couldn't publish mail")
		return models.TaskResponse{Success: false}, err
	}
	m.l.Debug().Msgf("mail pushed to MQ %+v", mail)

	return models.TaskResponse{Success: true}, nil
}

func (m *MQClient) PushTask(ctx context.Context, task *models.Task, action, kind int) (models.TaskResponse, error) {
	value := grpc_task.TaskMessage{
		TaskId:      task.ID.Hex(),
		Email:       task.Creator,
		Name:        task.Name,
		Description: task.Descr,
		CreateTs:    timestamppb.New(task.StartTime),
		Status:      grpc_task.Status(task.StatusType),
		EndTime:     timestamppb.New(task.EndTime),
		EmailList:   task.EmailList,
		Action:      grpc_task.Action(action),
		Kind:        grpc_task.Kind(kind),
	}
	valueBytes, err := proto.Marshal(&value)
	if err != nil {
		m.l.Warn().Err(err).Msg("MQClient.PushTask proto marshal problem")
		return models.TaskResponse{}, fmt.Errorf("couldn't encode mail in proto: %w", err)
	}
	msg := types.Message{
		Key:   []byte("push_task tasks"),
		Value: valueBytes,
	}
	if err := m.ProducerTaskTopic.SendMessage(ctx, []types.Message{msg}, partition); err != nil {
		m.l.Warn().Err(err).Msg("MQClient.PushTask couldn't publish task")
		return models.TaskResponse{Success: false}, err
	}
	m.l.Debug().Msgf("task pushed to MQ %+v", task)
	return models.TaskResponse{Success: true}, nil
}
