package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	e.Server.ReadHeaderTimeout = 5 * time.Second
	e.Server.ReadTimeout = 15 * time.Second
	e.Server.WriteTimeout = 15 * time.Second
	e.Server.IdleTimeout = 60 * time.Second
	e.Server.MaxHeaderBytes = 1 << 20
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
	svc := service.NewPlatformService(
		repo,
		adminEmails,
		os.Getenv("STRIPE_CHECKOUT_BASE_URL"),
		os.Getenv("STRIPE_SECRET_KEY"),
		os.Getenv("STRIPE_WEBHOOK_SECRET"),
		os.Getenv("STRIPE_SUCCESS_URL"),
		os.Getenv("STRIPE_CANCEL_URL"),
		os.Getenv("STRIPE_MEMBERSHIP_PRICE_ID"),
		os.Getenv("STRIPE_COURSE_PRICE_ID"),
	)
	ctl := controller.New(svc, os.Getenv("INTERNAL_INGEST_TOKEN"), os.Getenv("SUPABASE_JWT_SECRET"))

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":     "ok",
			"repository": repositoryMode,
		})
	})

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"service":    "those-forgotten-api",
			"status":     "ok",
			"repository": repositoryMode,
			"healthz":    "/healthz",
			"readyz":     "/readyz",
			"swagger":    "/swagger",
		})
	})

	e.GET("/swagger.yaml", func(c echo.Context) error {
		for _, candidate := range swaggerCandidates() {
			if _, err := os.Stat(candidate); err == nil {
				return c.File(candidate)
			}
		}

		return c.JSON(http.StatusNotFound, map[string]any{"error": "swagger spec not found"})
	})

	e.GET("/swagger", func(c echo.Context) error {
		return c.HTML(http.StatusOK, swaggerHTML())
	})

	e.GET("/readyz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"status":                  "ready",
			"repository":              repositoryMode,
			"jwtAuthConfigured":       os.Getenv("SUPABASE_JWT_SECRET") != "",
			"stripeWebhookConfigured": os.Getenv("STRIPE_WEBHOOK_SECRET") != "",
			"notificationWorkerMode":  "queue",
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

func swaggerCandidates() []string {
	executable, _ := os.Executable()

	return []string{
		strings.TrimSpace(os.Getenv("SWAGGER_PATH")),
		"docs/swagger.yaml",
		filepath.Join(filepath.Dir(executable), "docs", "swagger.yaml"),
		filepath.Join("services", "api-go", "docs", "swagger.yaml"),
	}
}

func swaggerHTML() string {
	return fmt.Sprintf(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Customs Auction Platform Swagger</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: "/swagger.yaml",
        dom_id: "#swagger-ui",
        deepLinking: true,
        presets: [SwaggerUIBundle.presets.apis]
      });
    </script>
  </body>
</html>`)
}
