package route

import (
	"github.com/labstack/echo/v4"
	"github.com/raflyritonga/file-storage-api/internal/handler"
	mw "github.com/raflyritonga/file-storage-api/internal/middleware"
	"github.com/raflyritonga/file-storage-api/internal/repository"
	"github.com/raflyritonga/file-storage-api/pkg/helper/auth"
)

func RegisterAuthRoutes(
	e *echo.Echo, 
	authHandler *handler.AuthHandler, 
	tokenMgr *auth.TokenManager, 
	sessionRepo repository.SessionRepository) {

	e.POST("/auth/register", authHandler.Register)
    e.POST("/auth/login", authHandler.Login)
    e.POST("/auth/refresh", authHandler.Refresh)

    protected := e.Group("")
    protected.Use(mw.RequireAuth(tokenMgr, sessionRepo))

    protected.POST("/auth/logout", authHandler.Logout)
    protected.GET("/auth/me", authHandler.Me)

    admin := protected.Group("/admin")
    admin.Use(mw.RequireRoles("admin"))
    admin.GET("/ping", func(c echo.Context) error {
        return c.JSON(200, map[string]string{"message": "admin ok"})
    })

	}