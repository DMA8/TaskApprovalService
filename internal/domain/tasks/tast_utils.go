package tasks

import (
	"context"
	"gitlab.com/g6834/team31/tasks/internal/domain/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	approve     = "approve"
	decline     = "decline"
	decisionURL = "/decision"
)

//TODO
func validateTasks(task *models.Task) error {
	if task.Creator == "" || task.Name == "" {
		return ErrTaskValidationError
	}
	return nil
}

type EmailMsg struct {
	TaskID      primitive.ObjectID `json:"taskID"`
	LinkApprove string             `json:"approve_link"`
	LinkDecline string             `json:"decline_link"`
}

func createUsers(task *models.Task) {
	for _, email := range task.EmailList {
		task.Users = append(task.Users, models.User{Email: email})
	}
}

// проверка прав на исправление задачи
func (s *Service) CheckAlterPermission(ctx context.Context, task *models.Task) error {
	oldTask, err := s.ReadTask(ctx, task)
	if err != nil {
		return err
	}
	if oldTask.Creator != task.Creator {
		return ErrNotTaskCreator
	}
	return nil
}
