package ports

import (
	"context"
	"gitlab.com/g6834/team31/tasks/internal/domain/models"
)

// ClientTask Интерфейс grpc Клиента
type ClientTask interface {
	PushTask(ctx context.Context, task *models.Task, action, kind int) (models.TaskResponse, error)
	PushMail(ctx context.Context, mail *models.Mail) (models.TaskResponse, error)
}
