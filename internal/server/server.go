package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ashtishad/xpay/internal/common"
	"github.com/ashtishad/xpay/internal/infra/postgres"
	"github.com/ashtishad/xpay/internal/server/middlewares"
	"github.com/ashtishad/xpay/internal/server/routes"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router     *gin.Engine
	httpServer *http.Server
	DB         *sql.DB
	Config     *common.AppConfig
	Logger     *slog.Logger
}

func NewServer(ctx context.Context) (*Server, error) {
	cfg, err := common.LoadConfig()
	if err != nil {
		return nil, err
	}

	gin.SetMode(cfg.App.GinMode)
	router := gin.New()

	logLevel := slog.LevelInfo

	if gin.IsDebugging() {
		logLevel = slog.LevelDebug
	}

	logger := common.NewSlogger(logLevel)
	slog.SetDefault(logger)

	db, err := postgres.NewConnection(ctx, cfg.DB, logger)
	if err != nil {
		return nil, err
	}

	if err := postgres.RunMigrations(ctx, db, logger); err != nil {
		logger.Warn("failed to run migrations", "err", err)
		return nil, err
	}

	_ = router.SetTrustedProxies(nil)

	s := &Server{
		Router: router,
		Logger: logger,
		DB:     db,
		Config: cfg,
		httpServer: &http.Server{
			Addr:         cfg.App.ServerAddress,
			Handler:      router,
			IdleTimeout:  common.Timeouts.Server.Read * 2,
			ReadTimeout:  common.Timeouts.Server.Read,
			WriteTimeout: common.Timeouts.Server.Write,
		},
	}

	s.setupMiddlewares()

	s.setupRoutes()

	s.Logger.Info(fmt.Sprintf("Swagger Specs available at %s/swagger/index.html", s.httpServer.Addr))

	return s, nil
}

func (s *Server) setupMiddlewares() {
	s.Router.Use(middlewares.InitMiddlewares(s.Logger)...)
}

func (s *Server) setupRoutes() {
	apiGroup := s.Router.Group("/api/v1")
	routes.InitRoutes(apiGroup, s.Logger, s.DB, s.Config)
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.DB.Close(); err != nil {
		s.Logger.Error("failed to close database connection", "error", err)
	}

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}
