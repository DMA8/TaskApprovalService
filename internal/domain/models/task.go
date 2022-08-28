package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Desicion int

const (
	Approve    Desicion = 1
	Decline    Desicion = -1
	NoDecision Desicion = 0
)

// описываем структуру таски
type Task struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Creator       string             `json:"creator" bson:"creator"`
	Name          string             `json:"name" bson:"name"`
	Descr         string             `json:"description,omitempty" bson:"description,omitempty"`
	StartTime     time.Time          `json:"ts_created,omitempty"`
	EndTime       time.Time          `json:"ts_finished,omitempty"`
	StatusType    int                `json:"status,omitempty" bson:"status,omitempty"`
	EmailList     []string           `json:"email_list,omitempty"`
	EmailProgress int                `json:"emailProgress,omitempty" bson:"emailProgress,omitempty"`
	Users         []User             `json:"users,omitempty" bson:"users,omitempty"`
}

type User struct {
	Email  string   `json:"email,omitempty" bson:"email,omitempty"`
	Status Desicion `json:"status,omitempty" bson:"status,omitempty"`
}

type TaskAnalytic struct {
	CTX context.Context
	Task   *Task
	Action int
	Kind   int
}

type TaskResponse struct {
	Success bool `json:"success"`
}

type Mail struct {
	CTX context.Context
	Header    string    `json:"header"`
	Body      string    `json:"body"`
	CreateTS  time.Time `json:"create_ts"`
	EmailList []string  `json:"email_list"`
}

type EmailMsg struct {
	TaskID          primitive.ObjectID `json:"taskID"`
	TaskName        string             `json:"task_name"`
	TaskDescription string             `json:"task_description"`
	LinkApprove     string             `json:"approve_link"`
	LinkDecline     string             `json:"decline_link"`
}
