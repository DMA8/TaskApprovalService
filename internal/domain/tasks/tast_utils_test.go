package tasks

import (
	"context"
	"testing"

	"gitlab.com/g6834/team31/tasks/internal/config"
	"gitlab.com/g6834/team31/tasks/internal/domain/models"
	mock_ports "gitlab.com/g6834/team31/tasks/internal/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceCreateTask(t *testing.T) {
	ctr := gomock.NewController(t)
	db := mock_ports.NewMockTasksStorage(ctr)
	cfg := config.Config{HTTP: config.HTTPConfig{URI: "test1", APIVersion: "tasks/v1"}}
	service := New(db, cfg)
	ctx := context.Background()
	tsk1 := &models.Task{Name: "test1", Creator: "test"}
	ans1 := &models.Task{Name: "test1", Creator: "test"}
	db.EXPECT().ReadTask(gomock.Any(), tsk1).Return(ans1, nil).Times(1)
	err := service.CheckAlterPermission(ctx, tsk1)
	require.NoError(t, err)

	ans2 := &models.Task{Name: "test1", Creator: "test2"}
	db.EXPECT().ReadTask(gomock.Any(), tsk1).Return(ans2, nil).Times(1)
	err = service.CheckAlterPermission(ctx, tsk1)
	require.Error(t, err)
}

func Test_createUsers(t *testing.T) {
	type args struct {
		task *models.Task
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{task: &models.Task{EmailList: []string{"e1", "e2"}}},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createUsers(tt.args.task)
			assert.Equal(t, len(tests[i].args.task.EmailList), len(tests[i].args.task.Users))
		})
	}
}
