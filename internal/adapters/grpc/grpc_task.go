package grpc

import (
	"context"
	"gitlab.com/g6834/team31/tasks/internal/domain/models"
	"gitlab.com/g6834/team31/tasks/pkg/grpc_task"

	"google.golang.org/protobuf/types/known/timestamppb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TaskClient struct {
	client grpc_task.GrpcTaskClient
	conn   *grpc.ClientConn
}

func NewServer(ctx context.Context, host, port string) (*TaskClient, error) {
	connStr := host + port
	conn, err := grpc.DialContext(ctx, connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := grpc_task.NewGrpcTaskClient(conn)
	return &TaskClient{
		client: client,
		conn:   conn,
	}, nil
}

func (t *TaskClient) Stop() error {
	return t.conn.Close()
}

func (t *TaskClient) PushTask(ctx context.Context, task *models.Task, action, kind int) (models.TaskResponse, error) {
	response, err := t.client.PushTask(ctx, &grpc_task.TaskMessage{
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
	})
	if err != nil {
		return models.TaskResponse{Success: false,}, err
	}
	return models.TaskResponse{Success: response.Success,}, nil
}

func (t *TaskClient) PushMail(ctx context.Context, mail *models.Mail) (models.TaskResponse, error) {
	response, err := t.client.PushMail(ctx, &grpc_task.Mail{
		Header:    mail.Header,
		Body:      mail.Body,
		CreateTs:  timestamppb.New(mail.CreateTS),
		EmailList: mail.EmailList,
	})
	if err != nil {
		return models.TaskResponse{
			Success: false,
		}, err
	}
	return models.TaskResponse{
		Success: response.Success,
	}, nil
}
