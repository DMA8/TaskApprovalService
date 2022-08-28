package ports

import (
	"context"

	"gitlab.com/g6834/team31/tasks/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Интерфейс описывающий сервис Tasks, который будут дергать адаптеры
type Tasks interface {
	CreateTask(ctx context.Context, task *models.Task) (primitive.ObjectID, error)
	ReadTask(ctx context.Context, task *models.Task) (*models.Task, error)
	ReadTaskById(ctx context.Context, task *models.Task) (*models.Task, error)
	UpdateTask(ctx context.Context, task *models.Task) (*models.Task, error)
	DeleteTask(ctx context.Context, task *models.Task) error
	ListTask(ctx context.Context, task *models.Task) ([]*models.Task, error)
	UpdateApprovalStatus(ctx context.Context, task *models.Task, approvalEmail string, decision models.Desicion) (*models.Task, error)
}
