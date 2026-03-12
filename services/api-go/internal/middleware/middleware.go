package middleware

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	echo "github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func Install(e *echo.Echo, origin string) {
	origins := []string{"http://localhost:3000"}
	if origin != "" {
		origins = []string{origin}
	}

	rateLimitRPS := 5
	if parsed, err := strconv.Atoi(strings.TrimSpace(os.Getenv("RATE_LIMIT_RPS"))); err == nil && parsed > 0 {
		rateLimitRPS = parsed
	}

	e.Use(echomw.Recover())
	e.Use(echomw.Logger())
	e.Use(echomw.RequestID())
	e.Use(echomw.Secure())
	e.Use(echomw.BodyLimit("1M"))
	e.Use(echomw.TimeoutWithConfig(echomw.TimeoutConfig{Timeout: 15 * time.Second}))
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
	}))
	e.Use(csrfOriginGuard(parseOrigins(origins)))
	e.Use(echomw.RateLimiter(echomw.NewRateLimiterMemoryStore(rate.Limit(rateLimitRPS))))
}

func csrfOriginGuard(trustedOrigins []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			method := c.Request().Method
			if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
				return next(c)
			}

			if c.Request().Header.Get("Cookie") == "" {
				return next(c)
			}

			origin := strings.TrimSpace(c.Request().Header.Get("Origin"))
			referer := strings.TrimSpace(c.Request().Header.Get("Referer"))
			if isTrustedOrigin(origin, trustedOrigins) || isTrustedOrigin(referer, trustedOrigins) {
				return next(c)
			}

			return echo.NewHTTPError(http.StatusForbidden, "csrf origin rejected")
		}
	}
}

func parseOrigins(defaultOrigins []string) []string {
	raw := strings.TrimSpace(os.Getenv("CSRF_TRUSTED_ORIGINS"))
	if raw == "" {
		return defaultOrigins
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		normalized := strings.TrimSpace(part)
		if normalized != "" {
			origins = append(origins, normalized)
		}
	}

	if len(origins) == 0 {
		return defaultOrigins
	}

	return origins
}

func isTrustedOrigin(candidate string, trustedOrigins []string) bool {
	if candidate == "" {
		return false
	}

	parsed, err := url.Parse(candidate)
	if err != nil {
		return false
	}

	normalized := parsed.Scheme + "://" + parsed.Host
	for _, trustedOrigin := range trustedOrigins {
		if normalized == strings.TrimRight(trustedOrigin, "/") {
			return true
		}
	}

	return false
}
