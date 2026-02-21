package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/raflyritonga/file-storage-api/internal/config"
	"github.com/raflyritonga/file-storage-api/internal/repository"

	"github.com/raflyritonga/file-storage-api/internal/handler"
	"github.com/raflyritonga/file-storage-api/internal/route"
	"github.com/raflyritonga/file-storage-api/internal/service"
	"github.com/raflyritonga/file-storage-api/pkg/helper"
	"github.com/raflyritonga/file-storage-api/pkg/helper/auth"
	redis "github.com/redis/go-redis/v9"
)

type Server struct {
	echo *echo.Echo
	cfg  config.Config
	db *sql.DB
	redis *redis.Client
}

func New(cfg config.Config) (*Server, error) {
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		code := http.StatusInternalServerError
		message := "internal server error"

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			switch value := he.Message.(type) {
			case string:
				if value != "" {
					message = value
				}
			case error:
				message = value.Error()
			default:
				message = fmt.Sprint(value)
			}
		}

		if jsonErr := c.JSON(code, helper.ErrorResponse(message)); jsonErr != nil {
			c.Logger().Error(jsonErr)
		}
	}
	e.Use(middleware.RequestID(), middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339}","request_id":"${id}","method":"${method}","uri":"${uri}","status":${status},"error":"${error}"}` + "\n",
		Skipper: func(c echo.Context) bool {
			path := c.Path()
			urlPath := c.Request().URL.Path
			return path == "/favicon.ico" || urlPath == "/favicon.ico"
		},
	}))

	db, err := config.NewDB(cfg)
	if err != nil {
		err := errors.New("failed to connect to postgres: " + err.Error())
		return nil, err
	}
	log.Println("Connected to PostgreSQL")

	redisClient, err := config.NewRedisClient(cfg)
	if err != nil {
		err := errors.New("failed to connect to redis: " + err.Error())
		return nil, err
	}
	log.Println("Connected to Redis")

	userRepo := repository.NewUserRepository(db)
    sessionRepo := repository.NewSessionRepository(redisClient)
	tokenMgr := auth.NewTokenManager(
        cfg.JWTAccessSecret,
        cfg.JWTRefreshSecret,
        cfg.JWTIssuer,
        cfg.JWTAccessTTLMin,
        cfg.JWTRefreshTTLHr,
    )

	authService := service.NewAuthService(userRepo, sessionRepo, tokenMgr)
    authHandler := handler.NewAuthHandler(authService, sessionRepo)

    route.RegisterAuthRoutes(e, authHandler, tokenMgr, sessionRepo)
	route.RegisterHealthRoutes(e)

    return &Server{echo: e, cfg: cfg, db: db, redis: redisClient}, nil
}

func (s *Server) Start() error {
	addr := "localhost:" + s.cfg.AppPort
	return s.echo.Start(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("Gracefully shutting down server...")
	return s.echo.Shutdown(ctx)
}
