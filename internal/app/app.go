package app

import (
	stdlog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"technical_test/config"
	v1 "technical_test/internal/controller/http/v1"
	"technical_test/internal/usecase"
	"technical_test/internal/usecase/repo"
	"technical_test/pkg/httpserver"
	"technical_test/pkg/logger"
	"technical_test/pkg/postgres"
)

func Run() {
	// Load environment variables
	_ = godotenv.Load()

	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		stdlog.Fatalf("Config error: %s", err)
	}

	// Logger
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, cfg.PG.PoolMax)
	if err != nil {
		l.Fatal("app - Run - postgres.New: %v", err)
	}
	defer pg.Close()

	// Usecase
	itemUseCase := usecase.NewItemUseCase(
		repo.NewItemRepo(pg),
	)
	companyUseCase := usecase.NewCompanyUseCase(
		repo.NewCompanyRepo(pg),
	)

	// HTTP Server
	handler := gin.Default()
	v1.NewRouter(handler, itemUseCase, companyUseCase, l)
	httpServer := httpserver.New(handler, cfg.HTTP.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error("app - Run - httpServer.Notify: %v", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error("app - Run - httpServer.Shutdown: %v", err)
	}
}
