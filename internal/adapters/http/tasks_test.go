package http_test

// import (
// 	"bytes"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// 	"tasks/internal/config"
// 	"tasks/internal/domain/models"
// 	mock_ports "tasks/internal/mocks"
// 	"tasks/internal/ports"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// )

// func NewServer(tasks ports.Tasks, config *config.Config) *Server {
// 	return &Server{
// 		tasksService: tasks,
// 		cfg:          config,
// 	}
// }

// func TestCreate(t *testing.T) {
// 	cfg := &config.Config{
// 		HTTP: config.HTTPConfig{
// 			URI: ":8800",
// 		},
// 	}
// 	ctr := gomock.NewController(t)
// 	mockTasks := mock_ports.NewMockTasks(ctr)
// 	newHandler := NewServer(mockTasks, cfg)
// 	handler := http.HandlerFunc(newHandler.create)
// 	recorder := httptest.NewRecorder()

// 	test := models.Task{
// 		Name:  "James",
// 		Descr: "Holiday for me, please",
// 	}
// 	// vals := url.Values{}
// 	// vals.Set("name", "James")
// 	// vals.Add("description", "The holiday for me, please")
// 	// reqBody := bytes.NewBufferString(vals.Encode())

// 	request, err := http.NewRequest("POST", fmt.Sprint("localhost:4000/tasks/v1/create"), reqBody)
// 	mockTasks.EXPECT().CreateTask()
// }
