package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"gitlab.com/g6834/team31/tasks/internal/config"
	"gitlab.com/g6834/team31/tasks/internal/domain/models"
	"gitlab.com/g6834/team31/auth/pkg/logging"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	IsError    bool   `json:"is_error"`
}

type MessageWithTask struct {
	StatusCode int            `json:"status_code"`
	Message    string         `json:"message"`
	Task       []*models.Task `json:"tasks"`
	IsError    bool           `json:"is_error"`
}

func initHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func WriteAnswer(writer http.ResponseWriter, status int, message string, l *logging.Logger) {
	var errorFlag bool
	if status >= 400 {
		errorFlag = true
	}
	msg := Message{
		StatusCode: status,
		Message:    message,
		IsError:    errorFlag,
	}
	writer.WriteHeader(status)
	err := json.NewEncoder(writer).Encode(msg)
	if err != nil {
		l.Warn().Err(err).Msgf("WriteAnswer err while encoding msg %+v", msg)
	}
}

func writeAnswerWithTask(writer http.ResponseWriter, status int, message string, task []*models.Task, l *logging.Logger) {
	var errorFlag bool
	if status >= 400 {
		errorFlag = true
	}
	msg := MessageWithTask{
		StatusCode: status,
		Message:    message,
		Task:       task,
		IsError:    errorFlag,
	}
	writer.WriteHeader(status)
	err := json.NewEncoder(writer).Encode(msg)
	if err != nil {
		l.Warn().Err(err).Msgf("writeAnswerWithTask err while encoding msg %+v", msg)
	}
}

func loginFromCtx(ctx context.Context, l *logging.Logger) (string, error) {
	usr := ctx.Value(userLogin("userLogin"))
	switch usr := usr.(type) {
	case userLogin, string:
		return usr.(string), nil
	default:
		l.Debug().Msgf("loginFromCtx: coudn't extract login from ctx %+v", ctx)
		return "", errors.New("coudn't get login from ctx")
	}
}

func decodeTaskFromCtx(reader io.Reader, logger *logging.Logger) (*models.Task, error) {
	var task models.Task
	err := json.NewDecoder(reader).Decode(&task)
	if err != nil {
		logger.Debug().Err(err).Msgf("decodeTaskFromCtx coudn't decode input %+v", reader)
		return nil, err
	}
	return &task, nil
}

func GenerateLink(taskID primitive.ObjectID, typeLink, email string, cfg config.HTTPConfig) string {
	id := taskID.Hex()
	return fmt.Sprintf("%s%s/%s/%s/%s", cfg.URI, cfg.APIVersion, id, typeLink, email)
}

var ErrFinishApprovement = fmt.Errorf("task is approved")
var ErrTaskDeclined = fmt.Errorf("task is declined")


// если все согласовали, возвращает -1
func NextEmailIndx(task *models.Task) (int, error) {
	for index, user := range task.Users {
		if user.Status == models.NoDecision {
			return index, nil
		} else if user.Status == models.Decline {
			return -1, ErrTaskDeclined
		}
	}
	return -1, ErrFinishApprovement
}
