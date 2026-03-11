package app

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/dennislee928/those_forgotten/services/api-go/internal/controller"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/middleware"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/repository"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/service"
	echo "github.com/labstack/echo/v4"
)

type Server struct {
	echo           *echo.Echo
	port           string
	cleanup        func()
	repositoryMode string
}

func NewServer() (*Server, error) {
	e := echo.New()
	middleware.Install(e, os.Getenv("WEB_ORIGIN"))

	repo := repository.Repository(repository.NewMemoryRepository())
	cleanup := func() {}
	repositoryMode := "memory"
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		postgresRepo, err := repository.NewPostgresRepository(context.Background(), databaseURL)
		if err == nil {
			repo = postgresRepo
			cleanup = postgresRepo.Close
			repositoryMode = "postgres"
		}
	}

	adminEmails := parseAdminEmails(os.Getenv("ADMIN_EMAILS"))
	svc := service.NewPlatformService(repo, adminEmails)
	ctl := controller.New(svc, os.Getenv("INTERNAL_INGEST_TOKEN"))

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":     "ok",
			"repository": repositoryMode,
		})
	})

	controller.RegisterRoutes(e, ctl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Server{echo: e, port: port, cleanup: cleanup, repositoryMode: repositoryMode}, nil
}

func parseAdminEmails(raw string) []string {
	if raw == "" {
		return nil
	}

	values := strings.Split(raw, ",")
	emails := make([]string, 0, len(values))
	for _, value := range values {
		normalized := strings.TrimSpace(strings.ToLower(value))
		if normalized != "" {
			emails = append(emails, normalized)
		}
	}

	return emails
}

func (s *Server) Start() error {
	defer s.cleanup()
	return s.echo.Start(":" + s.port)
}

func (s *Server) Echo() *echo.Echo {
	return s.echo
}
