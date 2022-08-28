// Code generated by MockGen. DO NOT EDIT.
// Source: internal/ports/tasks.go

// Package mock_ports is a generated GoMock package.
package mock_ports

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "gitlab.com/g6834/team31/tasks/internal/domain/models"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
)

// MockTasks is a mock of Tasks interface.
type MockTasks struct {
	ctrl     *gomock.Controller
	recorder *MockTasksMockRecorder
}

// MockTasksMockRecorder is the mock recorder for MockTasks.
type MockTasksMockRecorder struct {
	mock *MockTasks
}

// NewMockTasks creates a new mock instance.
func NewMockTasks(ctrl *gomock.Controller) *MockTasks {
	mock := &MockTasks{ctrl: ctrl}
	mock.recorder = &MockTasksMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTasks) EXPECT() *MockTasksMockRecorder {
	return m.recorder
}

// CreateTask mocks base method.
func (m *MockTasks) CreateTask(ctx context.Context, task *models.Task) (primitive.ObjectID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", ctx, task)
	ret0, _ := ret[0].(primitive.ObjectID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockTasksMockRecorder) CreateTask(ctx, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockTasks)(nil).CreateTask), ctx, task)
}

// DeleteTask mocks base method.
func (m *MockTasks) DeleteTask(ctx context.Context, task *models.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTask", ctx, task)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockTasksMockRecorder) DeleteTask(ctx, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockTasks)(nil).DeleteTask), ctx, task)
}

// ListTask mocks base method.
func (m *MockTasks) ListTask(ctx context.Context, task *models.Task) ([]*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTask", ctx, task)
	ret0, _ := ret[0].([]*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTask indicates an expected call of ListTask.
func (mr *MockTasksMockRecorder) ListTask(ctx, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTask", reflect.TypeOf((*MockTasks)(nil).ListTask), ctx, task)
}

// ReadTask mocks base method.
func (m *MockTasks) ReadTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadTask", ctx, task)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadTask indicates an expected call of ReadTask.
func (mr *MockTasksMockRecorder) ReadTask(ctx, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadTask", reflect.TypeOf((*MockTasks)(nil).ReadTask), ctx, task)
}

// ReadTaskById mocks base method.
func (m *MockTasks) ReadTaskById(ctx context.Context, task *models.Task) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadTaskById", ctx, task)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadTaskById indicates an expected call of ReadTaskById.
func (mr *MockTasksMockRecorder) ReadTaskById(ctx, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadTaskById", reflect.TypeOf((*MockTasks)(nil).ReadTaskById), ctx, task)
}

// UpdateApprovalStatus mocks base method.
func (m *MockTasks) UpdateApprovalStatus(ctx context.Context, task *models.Task, approvalEmail string, decision models.Desicion) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateApprovalStatus", ctx, task, approvalEmail, decision)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateApprovalStatus indicates an expected call of UpdateApprovalStatus.
func (mr *MockTasksMockRecorder) UpdateApprovalStatus(ctx, task, approvalEmail, decision interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateApprovalStatus", reflect.TypeOf((*MockTasks)(nil).UpdateApprovalStatus), ctx, task, approvalEmail, decision)
}

// UpdateTask mocks base method.
func (m *MockTasks) UpdateTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", ctx, task)
	ret0, _ := ret[0].(*models.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockTasksMockRecorder) UpdateTask(ctx, task interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockTasks)(nil).UpdateTask), ctx, task)
}