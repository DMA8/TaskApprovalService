package tasks

import (
	"context"
	"gitlab.com/g6834/team31/tasks/internal/config"
	"gitlab.com/g6834/team31/tasks/internal/domain/models"
	mock_ports "gitlab.com/g6834/team31/tasks/internal/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateTask(t *testing.T) {
	ctr := gomock.NewController(t)
	db := mock_ports.NewMockTasksStorage(ctr)
	cfg := config.Config{HTTP: config.HTTPConfig{URI: "test1", APIVersion: "tasks/v1"}}
	service := New(db, cfg)
	ctx := context.Background()
	tsk1 := &models.Task{Name: "test1", Creator: "test1"}
	ans1 := primitive.NewObjectID()
	db.EXPECT().CreateTask(gomock.Any(), tsk1).Return(ans1, nil).Times(1)
	_, err := service.CreateTask(ctx, tsk1)
	require.NoError(t, err)

	task2 := &models.Task{Name: "test1", Creator: ""}
	_, err = service.CreateTask(ctx, task2)
	require.Error(t, err)
}
