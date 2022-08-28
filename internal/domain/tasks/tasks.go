package tasks

import (
	"context"
	"errors"

	"gitlab.com/g6834/team31/tasks/internal/config"
	"gitlab.com/g6834/team31/tasks/internal/domain/models"
	"gitlab.com/g6834/team31/tasks/internal/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrTaskIsAlreadyExistsInStorage = errors.New("such task is already in storage")
	ErrNotTaskCreator               = errors.New("task can be changed only by its creator")
	ErrTaskStatusIsAlreadySet       = errors.New("task status is already updated")
	ErrTaskValidationError          = errors.New("task should have name and its creator")
)

type Service struct {
	cfg config.Config
	db  ports.TasksStorage
}

func New(db ports.TasksStorage, cfg config.Config) *Service {
	return &Service{
		db:  db,
		cfg: cfg,
	}
}

//Каждому участнику при создании отправить 2 уникальные ссылки
func (s *Service) CreateTask(ctx context.Context, task *models.Task) (primitive.ObjectID, error) {
	ctx, span := otel.Tracer("team31_tasks").Start(ctx, "service.CreateTask")
	defer span.End()
	if err := validateTasks(task); err != nil {
		return primitive.ObjectID{}, err
	}
	createUsers(task)
	DBTaskID, err := s.db.CreateTask(ctx, task) // сначала кидаем в базу, и только потом отправляем
	if err != nil {
		return primitive.ObjectID{}, err
	}
	task.ID = DBTaskID
	span.SetAttributes(attribute.KeyValue{Key: "task_id", Value: attribute.StringValue(DBTaskID.String())})

	return DBTaskID, err
}

func (s *Service) ReadTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	return s.db.ReadTask(ctx, task)
}

func (s *Service) ReadTaskById(ctx context.Context, task *models.Task) (*models.Task, error) {
	return s.db.ReadTaskById(ctx, task)
}

// UpdateTask only for author
func (s *Service) UpdateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	err := s.CheckAlterPermission(ctx, task)
	if err != nil {
		return nil, err
	}
	return s.db.UpdateTask(ctx, task)
}

// DeleteTask only for author
func (s *Service) DeleteTask(ctx context.Context, task *models.Task) error {
	err := s.CheckAlterPermission(ctx, task)
	if err != nil {
		return err
	}
	return s.db.DeleteTask(ctx, task)
}

func (s *Service) ListTask(ctx context.Context, task *models.Task) ([]*models.Task, error) {
	return s.db.ListTask(ctx, task)
}

// updates переписать с $set, чтобы не ходить перед обновлением в базу
// добавить в tasks возврат следующего получателя письма
// проверять id отправителя и id в ссылке
func (s *Service) UpdateApprovalStatus(ctx context.Context, task *models.Task, approvalEmail string, decision models.Desicion) (*models.Task, error) {
	oldTask, err := s.ReadTaskById(ctx, task)
	if err != nil {
		return nil, err
	}
	for userInd, user := range oldTask.Users {
		if user.Email == approvalEmail {
			if user.Status != 0 {
				return nil, ErrTaskStatusIsAlreadySet
			}
			oldTask.Users[userInd].Status = decision
		}
	}
	return s.UpdateTask(ctx, oldTask)
}
