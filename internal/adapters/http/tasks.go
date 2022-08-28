package http

import (
	"fmt"
	"net/http"
	"time"

	"gitlab.com/g6834/team31/tasks/internal/config"
	"gitlab.com/g6834/team31/tasks/internal/domain/models"
	"go.opentelemetry.io/otel/attribute"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Server) tasksHandlers() http.Handler {
	h := chi.NewMux()
	h.Use(s.ValidateToken)
	h.Route("/", func(r chi.Router) {
		h.Post("/task", s.create)
		h.Get("/task/{taskName}", s.read)
		h.Put("/task", s.update)
		h.Delete("/task/{taskName}", s.delete)
		h.Get("/tasks", s.list)
		h.Get("/task/{taskID}/approve/{approvalLogin}", s.approve)
		h.Get("/task/{taskID}/decline/{approvalLogin}", s.decline)
	})
	return h
}

// create godoc
// @ID create
// @tags tasks
// @Summary create new task
// @Description create new task and send mail to approvals
// @Body {object} models.Task
// @Success 200 {object} Message
// @Failure 403 {object} Message
// @Failure 500 {object} Message
// @Router /task [post]
// @Accept       json
// @Produce      json
// @Param input body models.Task true "account info"
func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "server.create handler")
	defer span.End()
	initHeaders(w)
	inpTask, err := decodeTaskFromCtx(r.Body, s.logger)
	if err != nil {
		WriteAnswer(w, http.StatusBadRequest, fmt.Sprintf("couldn't parse query:%s", err), s.logger)
		return
	}
	inpTask.Creator, err = loginFromCtx(ctx, s.logger)
	span.SetAttributes(attribute.KeyValue{Key: "user_login", Value: attribute.StringValue(inpTask.Creator)})

	if err != nil {
		WriteAnswer(w, http.StatusBadRequest, fmt.Sprintf("couldn't get userlogin from ctx:%s", err), s.logger)
		return
	}
	inpTask.StartTime = time.Now()
	inpTask.ID, err = s.tasksService.CreateTask(ctx, inpTask)
	if err != nil {
		s.logger.Debug().Msgf("s.create couldn't create task %s", err.Error())
		WriteAnswer(w, http.StatusInternalServerError, fmt.Sprintf("couldn't create task! service error:%s", err), s.logger)
		return
	}
	WriteAnswer(w, http.StatusOK, "task is created successfully!", s.logger)
	go func() {
		s.OutGateway.Tasks <- &models.TaskAnalytic{CTX:ctx, Task: inpTask, Kind: 0, Action: 0}
	}()
	go func() {
		var firstEmailIndex int
		mailMessage, err := generateMailMessage(inpTask, firstEmailIndex, s.cfg)
		mailMessage.CTX = ctx
		if err != nil {
			s.logger.Warn().Msgf("s.create bad email index")
			return
		}
		s.logger.Debug().Msgf("mailMessage generated! %+v", *mailMessage)
		s.OutGateway.Mails <- mailMessage
	}()

}

func generateMailMessage(task *models.Task, emailIndex int, cfg config.HTTPConfig) (*models.Mail, error) {
	if emailIndex < 0 || emailIndex >= len(task.EmailList) {
		return nil, fmt.Errorf("wrong email index")
	}
	approveLink := GenerateLink(task.ID, "approve", task.EmailList[emailIndex], cfg)
	declineLink := GenerateLink(task.ID, "decline", task.EmailList[emailIndex], cfg)
	return &models.Mail{
		Header:    fmt.Sprintf("task name: %s taskID: %s", task.Name, task.ID.String()),
		Body:      fmt.Sprintf("your order %d, description: %s, to approve press: %s ; to decline press :%s ", emailIndex, task.Descr, approveLink, declineLink),
		CreateTS:  task.StartTime,
		EmailList: task.EmailList,
	}, nil
}

// read godoc
// @ID read
// @tags tasks
// @Summary read task
// @Description read existing task
// @Success 200 {object} models.Task
// @Failure 403 {object} Message
// @Failure 500 {object} Message
// @Router /task/{taskName} [get]
// @Produce      json
// @Param taskName path string true "Task name"
func (s *Server) read(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	taskName := chi.URLParam(r, "taskName")
	if taskName == "" {
		s.logger.Debug().Msgf("s.read. Bad input: taskName is empty")
		WriteAnswer(w, http.StatusBadRequest, "couldn't find taskName", s.logger)
		return
	}
	login, err := loginFromCtx(r.Context(), s.logger)
	if err != nil {
		s.logger.Debug().Err(err).Msgf("s.read. couldn't determine user login")
		WriteAnswer(w, http.StatusBadRequest, "s.read. couldn't determine user login", s.logger)
		return
	}
	t := models.Task{
		Creator: login,
		Name:    taskName,
	}
	taskFromDB, err := s.tasksService.ReadTask(r.Context(), &t)
	if err != nil {
		s.logger.Debug().Err(err).Msgf("s.read couldn't read task")
		WriteAnswer(w, http.StatusInternalServerError, fmt.Sprintf("couldn't read task! service error:%s", err), s.logger)
		return
	}
	writeAnswerWithTask(w, http.StatusOK, "task found!", []*models.Task{taskFromDB}, s.logger)
}

// update godoc
// @ID update
// @tags tasks
// @Summary update task
// @Description update existing task
// @Body {object} models.Task
// @Success 200 {object} Message
// @Failure 403 {object} Message
// @Failure 500 {object} Message
// @Router /task [put]
// @Accept       json
// @Produce      json
// @Param input body models.Task true "account info"
func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	inpTask, err := decodeTaskFromCtx(r.Body, s.logger)
	if err != nil {
		WriteAnswer(w, http.StatusBadRequest, fmt.Sprintf("couldn't parse query:%s", err), s.logger)
		return
	}
	inpTask.Creator, err = loginFromCtx(r.Context(), s.logger)
	if err != nil {
		WriteAnswer(w, http.StatusInternalServerError, "some internal problems", s.logger)
		return
	}
	updatedTask, err := s.tasksService.UpdateTask(r.Context(), inpTask)
	if err != nil {
		s.logger.Debug().Err(err).Msgf("s.update couldn't update task %+v,  err:%s", inpTask)
		WriteAnswer(w, http.StatusInternalServerError, fmt.Sprintf("couldn't update task! service error:%s", err), s.logger)
		return
	}
	writeAnswerWithTask(w, http.StatusOK, "task is updated successfully!", []*models.Task{updatedTask}, s.logger)

	go func() {
		s.OutGateway.Tasks <- &models.TaskAnalytic{Task: inpTask, Kind: 0, Action: 1}
	}()
}

// delete godoc
// @ID delete
// @tags tasks
// @Summary delete task
// @Description delete existing task
// @Body {object} models.Task
// @Success 200 {object} Message
// @Failure 403 {object} Message
// @Failure 500 {object} Message
// @Router /task/{taskName} [delete]
// @Accept       json
// @Produce      json
// @Param taskName path string true "Task name"
func (s *Server) delete(w http.ResponseWriter, r *http.Request) {
	initHeaders(w)
	taskName := chi.URLParam(r, "taskName")
	if taskName == "" {
		s.logger.Debug().Msgf("s.delete. Bad input: taskName is empty")
		WriteAnswer(w, http.StatusBadRequest, "couldn't find taskName", s.logger)
		return
	}
	login, err := loginFromCtx(r.Context(), s.logger)
	if err != nil {
		s.logger.Debug().Err(err).Msgf("s.delete. couldn't determine user login")
		WriteAnswer(w, http.StatusBadRequest, "s.delete. couldn't determine user login", s.logger)
		return
	}
	t := models.Task{
		Creator: login,
		Name:    taskName,
	}
	err = s.tasksService.DeleteTask(r.Context(), &t)
	if err != nil {
		s.logger.Debug().Err(err).Msgf("s.delete couldn't delete task %+v")
		WriteAnswer(w, http.StatusInternalServerError, fmt.Sprintf("couldn't update task! service error:%s", err), s.logger)
		return
	}
	WriteAnswer(w, http.StatusOK, "task is deleted successfully!", s.logger)

	go func() {
		s.OutGateway.Tasks <- &models.TaskAnalytic{Task: &t, Kind: 0, Action: 2}
	}()
}

// list godoc
// @ID list
// @tags tasks
// @Summary Get task list created by user
// @Description Get list of existing tasks
// @Success 200 {array} models.Task
// @Failure 403 {object} Message
// @Failure 500 {object} Message
// @Router /tasks [get]
func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	initHeaders(w)
	login, err := loginFromCtx(r.Context(), s.logger)
	if err != nil {
		WriteAnswer(w, http.StatusInternalServerError, "some internal errors :(", s.logger)
		return
	}
	task.Creator = login
	tasks, err := s.tasksService.ListTask(r.Context(), &task)
	if err != nil {
		s.logger.Debug().Err(err).Msgf("s.list couldn't list tasks %+v")
		WriteAnswer(w, http.StatusInternalServerError, fmt.Sprintf("couldn't list tasks! service error:%s", err), s.logger)
		return
	}
	writeAnswerWithTask(w, http.StatusOK, "tasks is listed successfully!", tasks, s.logger)
}

//отправить событие на отправку письма следующему получателю, либо закрыть таску, если полностью согласована
// approve godoc
// @ID approve
// @tags tasks
// @Summary approve current task
// @Description approve current task by current user
// @Produce json
// @Success 200 {object} Message
// @Failure 403 {object} Message
// @Failure 500 {object} Message
// @Param taskID path string true "Task ID"
// @Param approvalLogin path string true "Approval Login"
// @Router /task/{taskID}/approve/{approvalLogin} [get]
func (s *Server) approve(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var err error
	initHeaders(w)
	taskID := chi.URLParam(r, "taskID")
	approvalLogin := chi.URLParam(r, "approvalLogin")
	if taskID == "" || approvalLogin == "" {
		s.logger.Debug().Msgf("s.approve. Bad input: taskID: '%s' approvalLogin: '%s'", taskID, approvalLogin)
		WriteAnswer(w, http.StatusBadRequest, fmt.Sprintf("couldn't get taskID(%s) or approvalLogin(%s)", taskID, approvalLogin), s.logger)
		return
	}
	task.ID, err = primitive.ObjectIDFromHex(taskID)
	if err != nil {
		s.logger.Debug().Err(err).Msgf("s.approve. couldn't create task.ID")
		WriteAnswer(w, http.StatusBadRequest, fmt.Sprintf("couldn't cast task.ID fromHexString approvalLogin :%s", err), s.logger)
		return
	}
	updatedTask, err := s.tasksService.UpdateApprovalStatus(r.Context(), &task, approvalLogin, models.Approve)
	if err != nil {
		s.logger.Debug().Err(err).Msgf("s.approve couldn't update approval status err")
		WriteAnswer(w, http.StatusInternalServerError, fmt.Sprintf("couldn't update status taskID:%s, approvalLogin:%s:%s", task.ID, approvalLogin, err), s.logger)
		return
	}
	WriteAnswer(w, http.StatusOK, "ok", s.logger)

	// Делаем для будущей унификации в кафке
	task.Creator = approvalLogin
	task.StatusType = 1
	task.StartTime = time.Now()
	go func() {
		s.OutGateway.Tasks <- &models.TaskAnalytic{Task: &task, Kind: 1, Action: 0}
	}()
	go func() {
		emailIndex, err := NextEmailIndx(updatedTask)
		if err == ErrFinishApprovement {
			s.logger.Debug().Err(err).Msgf("Task is finished %+v", updatedTask)
			return
		} else if err == ErrTaskDeclined {
			s.logger.Debug().Err(err).Msgf("Task was declined %+v", updatedTask)
			return
		} else if err != nil {
			s.logger.Warn().Err(err).Msgf("undescribed error %s", err.Error())
			return
		}
		mailMessage, err := generateMailMessage(&task, emailIndex, s.cfg)
		if err != nil {
			s.logger.Warn().Err(err).Msgf("s.create bad email index")
			return
		}
		s.logger.Debug().Msgf("mail is sent to user %+v", *mailMessage)
		s.OutGateway.Mails <- mailMessage
	}()
}

// разослать всем об неуспешности согласования
// decline godoc
// @ID decline
// @tags tasks
// @Summary decline current task
// @Description decline current task by current user
// @Produce json
// @Success 200 {object} Message
// @Failure 403 {object} Message
// @Failure 500 {object} Message
// @Param taskID path string true "Task ID"
// @Param approvalLogin path string true "Approval Login"
// @Router /task/{taskID}/decline/{approvalLogin} [get]
func (s *Server) decline(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var err error
	initHeaders(w)
	taskID := chi.URLParam(r, "taskID")
	approvalLogin := chi.URLParam(r, "approvalLogin")
	if taskID == "" || approvalLogin == "" {
		s.logger.Debug().Msgf("s.decline. Bad input: taskID: '%s' approvalLogin: '%s'", taskID, approvalLogin)
		WriteAnswer(w, http.StatusBadRequest, fmt.Sprintf("coudn't get taskID or approvalLogin:%s", fmt.Errorf("err")), s.logger)
		return
	}
	task.ID, err = primitive.ObjectIDFromHex(taskID)
	if err != nil {
		s.logger.Debug().Err(err).Msgf("s.decline. coudn't create task.ID")
		WriteAnswer(w, http.StatusBadRequest, fmt.Sprintf("coudn't cast task.ID fromHexString approvalLogin :%s", err), s.logger)
		return
	}
	task.StatusType = -1
	_, err = s.tasksService.UpdateApprovalStatus(r.Context(), &task, approvalLogin, models.Decline)
	if err != nil {
		s.logger.Debug().Err(err).Msgf("s.decline coudn't update approval status err")
		WriteAnswer(w, http.StatusInternalServerError, fmt.Sprintf("couldn't update status taskID:%s, approvalLogin:%s: err: %v", task.ID, approvalLogin, err), s.logger)
		return
	}
	WriteAnswer(w, http.StatusOK, "ok", s.logger)

	// Делаем для будущей унификации в кафке
	task.Creator = approvalLogin
	task.StatusType = -1
	task.StartTime = time.Now()

	go func() {
		s.OutGateway.Tasks <- &models.TaskAnalytic{Task: &task, Kind: 1, Action: 0}
	}()
}
