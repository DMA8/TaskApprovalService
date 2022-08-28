package gateway

import (
	"context"
	"gitlab.com/g6834/team31/tasks/internal/domain/models"
	"gitlab.com/g6834/team31/tasks/internal/ports"
	"gitlab.com/g6834/team31/auth/pkg/logging"
)

// просто храним буф канал, в который кидаем события для аналитики

const (
	buffSize = 2
	nWorkers = 2
)

type Gateway struct {
	Tasks  chan *models.TaskAnalytic // ограничить только на запись?
	Mails  chan *models.Mail
	logger *logging.Logger
}

func New(l *logging.Logger) *Gateway {
	return &Gateway{
		Tasks:  make(chan *models.TaskAnalytic, buffSize),
		Mails:  make(chan *models.Mail, buffSize),
		logger: l,
	}
}

func (g *Gateway) StartGateway(ctx context.Context, p ports.ClientTask) {
	for i := 0; i < nWorkers; i++ {
		go g.SendWorker(ctx, p)
	}
	g.logger.Debug().Msgf("gateway StartGateway: launched %d send workers", nWorkers)
}

func (g *Gateway) SendWorker(ctx context.Context, p ports.ClientTask) {
	for {
		select {
		case task := <-g.Tasks:
			_, err := p.PushTask(ctx, task.Task, task.Action, task.Kind)
			if err != nil {
				g.logger.Warn().Err(err).Msgf("error while pushing task %+v", *task)
			} else {
				g.logger.Debug().Msgf("gateway.SendWorker task is sent away! %+v", *task)
			}
		case mail := <-g.Mails:
			_, err := p.PushMail(ctx, mail)
			if err != nil {
				g.logger.Warn().Err(err).Msgf("error while pushing task %+v", *mail,)
			} else {
				g.logger.Debug().Msgf("gateway.SendWorker mail is sent away! %+v", *mail)
			}
		case <-ctx.Done():
			g.logger.Info().Msgf("closing worker")
			return
		}
	}
}
