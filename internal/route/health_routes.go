package route

import (
	"github.com/labstack/echo/v4"
	"github.com/raflyritonga/file-storage-api/internal/handler"
)

func RegisterHealthRoutes(e *echo.Echo) {
	healthHandler := handler.NewHealthHandler()

	e.GET("/health", healthHandler.Check)
}