package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"gitlab.com/g6834/team31/auth/pkg/logging"
	_ "gitlab.com/g6834/team31/tasks/docs"
	"gitlab.com/g6834/team31/tasks/internal/adapters/gateway"
	"gitlab.com/g6834/team31/tasks/internal/config"
	"gitlab.com/g6834/team31/tasks/internal/ports"
)

type Server struct {
	tasksService ports.Tasks
	server       *http.Server
	Addr         string
	AuthClient   ports.ClientAuth
	OutGateway   *gateway.Gateway
	logger       *logging.Logger
	cfg          config.HTTPConfig
}

func New(cfg config.HTTPConfig, tasks ports.Tasks, client ports.ClientAuth, gateway *gateway.Gateway, log *logging.Logger) *Server {
	var s Server
	s.tasksService = tasks
	s.logger = log
	s.AuthClient = client
	s.cfg = cfg
	s.server = &http.Server{
		Handler: s.routes(),
		Addr:    cfg.URI,
	}
	s.OutGateway = gateway
	return &s
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) routes() http.Handler {
	r := chi.NewMux()
	cfg := config.NewConfig().HTTP
	r.Use(RequestID)
	r.Use(Logger(s.logger))
	r.Use(Prometheus())
	r.Use(TracingMiddleware)
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(":4000/swagger/doc.json")))
	r.Get("/healthz", s.HealthzHandler)
	r.Mount(cfg.APIVersion, s.tasksHandlers())
	return r
}

func (s *Server) HealthzHandler(w http.ResponseWriter, r *http.Request) {
	WriteAnswer(w, http.StatusOK, "жив", s.logger)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Start(ctx context.Context) error {
	return s.server.ListenAndServe()
}
