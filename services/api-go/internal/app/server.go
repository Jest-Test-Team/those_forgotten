package app

import (
	"context"
	"net/http"
	"os"

	"github.com/dennislee928/those_forgotten/services/api-go/internal/controller"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/middleware"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/repository"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/service"
	echo "github.com/labstack/echo/v4"
)

type Server struct {
	echo    *echo.Echo
	port    string
	cleanup func()
}

func NewServer() (*Server, error) {
	e := echo.New()
	middleware.Install(e, os.Getenv("WEB_ORIGIN"))

	repo := repository.Repository(repository.NewMemoryRepository())
	cleanup := func() {}
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		postgresRepo, err := repository.NewPostgresRepository(context.Background(), databaseURL)
		if err == nil {
			repo = postgresRepo
			cleanup = postgresRepo.Close
		}
	}

	svc := service.NewPlatformService(repo)
	ctl := controller.New(svc, os.Getenv("INTERNAL_INGEST_TOKEN"))

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	controller.RegisterRoutes(e, ctl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Server{echo: e, port: port, cleanup: cleanup}, nil
}

func (s *Server) Start() error {
	defer s.cleanup()
	return s.echo.Start(":" + s.port)
}

func (s *Server) Echo() *echo.Echo {
	return s.echo
}
