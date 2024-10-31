package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RateLimiter() echo.MiddlewareFunc {
	limiterStore := middleware.NewRateLimiterMemoryStore(20)
	return middleware.RateLimiter(limiterStore)
}
