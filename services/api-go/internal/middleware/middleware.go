package middleware

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

func Install(e *echo.Echo, origin string) {
	origins := []string{"http://localhost:3000"}
	if origin != "" {
		origins = []string{origin}
	}

	e.Use(echomw.Recover())
	e.Use(echomw.Logger())
	e.Use(echomw.Secure())
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: origins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
	}))
	e.Use(echomw.RateLimiter(echomw.NewRateLimiterMemoryStore(5)))
}
