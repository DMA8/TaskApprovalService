package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/g6834/team31/tasks/internal/domain/models"
	"gitlab.com/g6834/team31/auth/pkg/logging"
)

func TestWriteAnswer(t *testing.T) {
	reqRec := &httptest.ResponseRecorder{}
	testStatus := http.StatusOK
	testString := "hello"
	l := logging.New("debug")
	WriteAnswer(reqRec, testStatus, testString, &l)
	assert.Equal(t, testStatus, reqRec.Result().StatusCode)
}

func TestWriteanswerWithTask(t *testing.T) {
	task1 := &models.Task{
		Name:  "James",
		Descr: "Biba",
	}
	task2 := &models.Task{
		Name:  "Kirk",
		Descr: "Boba",
	}
	a := []*models.Task{}
	a = append(a, task1, task2)
	reqRec := &httptest.ResponseRecorder{}
	testStatus := 200
	testString := "hello"
	l := logging.New("debug")
	writeAnswerWithTask(reqRec, testStatus, testString, a, &l)
	assert.Equal(t, testStatus, reqRec.Result().StatusCode)
}

func TestLoginFromCtx(t *testing.T) {
	str := "test"
	ul := userLogin("userLogin")
	ctx := context.Background()
	l := logging.New("debug")

	ctx = context.WithValue(ctx, ul, str)
	testStr, err := loginFromCtx(ctx, &l)
	assert.Equal(t, str, testStr)
	assert.NoError(t, err)

	num := 123
	ctx = context.WithValue(ctx, ul, num)
	testStr2, err2 := loginFromCtx(ctx, &l)
	assert.Equal(t, "", testStr2)
	assert.Error(t, err2)
}
