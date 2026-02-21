package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/raflyritonga/file-storage-api/internal/dto"
	"github.com/raflyritonga/file-storage-api/internal/repository"
	"github.com/raflyritonga/file-storage-api/internal/service"
	"github.com/raflyritonga/file-storage-api/pkg/helper"
)

type AuthHandler struct {
	authService *service.AuthService
	sessionRepo repository.SessionRepository
}

func NewAuthHandler(authService *service.AuthService, sessionRepo repository.SessionRepository) *AuthHandler {
	return &AuthHandler{authService: authService, sessionRepo: sessionRepo}
}

func (h *AuthHandler) Register(ctx echo.Context) error {
	var req dto.RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid request"})
	}
	if err := h.authService.Register(ctx.Request().Context(), req); err != nil {
		return ctx.JSON(
			400, 
			helper.ErrorResponse(err.Error())		)
	}
	
	return ctx.JSON(201, helper.SuccessResponse("registration successful", nil))
}

func (h *AuthHandler) Login(ctx echo.Context) error {
	var req dto.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid payload"})
	}
	tokens, err := h.authService.Login(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(401, helper.ErrorResponse(err.Error()))
	}
	return ctx.JSON(200, helper.SuccessResponse("login successful", tokens))
}

func (h *AuthHandler) Logout(ctx echo.Context) error {
	accessJTI, _ := ctx.Get("access_jti").(string)
	accessExp, _ := ctx.Get("access_exp").(time.Time)
	ttl := int64(time.Until(accessExp).Seconds())

	if ttl > 0 {
		_ = h.sessionRepo.BlacklistAccessJTI(ctx.Request().Context(), accessJTI, ttl)
	}

	return ctx.JSON(200, helper.SuccessResponse("logout successful", nil))
}

func (h *AuthHandler) Refresh(ctx echo.Context) error {
	var req dto.RefreshRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, map[string]string{"error": "invalid payload"})
	}
	tokens, err := h.authService.Refresh(ctx.Request().Context(), req.RefreshToken)
	if err != nil {
		return ctx.JSON(401, helper.ErrorResponse(err.Error()))
	}
	return ctx.JSON(200, helper.SuccessResponse("token refreshed", tokens))
}

func (h *AuthHandler) Me(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]any{
		"user_id": ctx.Get("user_id"),
		"role":    ctx.Get("role"),
	})
}
