package application

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"gitlab.com/g6834/team31/tasks/internal/adapters/gateway"
	"gitlab.com/g6834/team31/tasks/internal/adapters/grpc"
	adapterMQ "gitlab.com/g6834/team31/tasks/internal/adapters/mq"

	"github.com/getsentry/sentry-go"
	"gitlab.com/g6834/team31/auth/pkg/logging"
	"gitlab.com/g6834/team31/tasks/internal/adapters/http"
	"gitlab.com/g6834/team31/tasks/internal/adapters/mongodb"
	"gitlab.com/g6834/team31/tasks/internal/config"
	"gitlab.com/g6834/team31/tasks/internal/domain/tasks"
	"gitlab.com/g6834/team31/tasks/pkg/mq"
)

func Start(ctx context.Context) {
	var (
		db  *mongodb.Database
		err error
	)

	ctx, cancel := context.WithCancel(ctx)
	cfg := config.NewConfig()
	logger := logging.New(cfg.Log.Level)
	logger.Info().Msg("tasks is starting...")
	logger.Debug().Msgf("%+v", cfg)
	db, err = mongodb.New(ctx, cfg.Mongo)

	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	logger.Info().Msg("db is ok")

	err = sentry.Init(sentry.ClientOptions{
		Dsn:   "https://fc7ec598813842d2b333029ad6810e3b@sentry.k8s.golang-mts-teta.ru/35",
		Debug: true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(1 * time.Second)
	grpcClientAuth, err := grpc.New(ctx, cfg.AuthGRPC.Host, cfg.AuthGRPC.Port)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}
	logger.Info().Msg("grcpAuthClient is ok")
	
	//gateway для grpc/кафки
	gateway := gateway.New(&logger)
	logger.Info().Msg("grpsAnalitycClient is ok")
	if err != nil {
		log.Fatal(err)
	}
	clientMQ, err := initClientMQ(cfg.Kafka, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("couldn't init MQ client")
	}
	gateway.StartGateway(ctx, clientMQ)
	tasksService := tasks.New(db, *cfg)
	httpServer := http.New(cfg.HTTP, tasksService, grpcClientAuth, gateway, &logger)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		logger.Info().Msgf("shutting down the app")
		cancel()
		err := httpServer.Shutdown(ctx)
		if err != nil {
			logger.Fatal().Msg(err.Error())
		}
		os.Exit(1)
	}()
	logger.Info().Msgf("launching http server")
	if err := httpServer.Start(ctx); err != nil {
		logger.Fatal().Msg(err.Error())
	}
}

func initClientMQ(cfg config.KafkaConfig, l *logging.Logger) (*adapterMQ.MQClient, error) {
	pMail, err := mq.NewProducer([]string{cfg.URL}, cfg.MailTopic)
	if err != nil {
		return nil, err
	}
	cMail, err := mq.NewConsumer([]string{cfg.URL}, cfg.MailTopic, cfg.GroupID)
	if err != nil {
		return nil, err
	}
	pTasks, err := mq.NewProducer([]string{cfg.URL}, cfg.TaskTopic)
	if err != nil {
		return nil, err
	}
	return adapterMQ.NewMQClient(cMail, pMail, pTasks, l), nil
}
