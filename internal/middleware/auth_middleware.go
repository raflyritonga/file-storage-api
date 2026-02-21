package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/raflyritonga/file-storage-api/internal/repository"
	"github.com/raflyritonga/file-storage-api/pkg/helper/auth"
)

func RequireAuth(tokenMgr *auth.TokenManager, sessionRepo repository.SessionRepository) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            authHeader := c.Request().Header.Get("Authorization")
            if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
            }

            parts := strings.SplitN(authHeader, " ", 2)
            if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header")
            }

            claims, err := tokenMgr.ParseAccessToken(parts[1])
            if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid access token")
            }

            blacklisted, err := sessionRepo.IsAccessJTIBlacklisted(c.Request().Context(), claims.ID)
            if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "cannot validate token")
            }
            if blacklisted {
				return echo.NewHTTPError(http.StatusUnauthorized, "token revoked")
            }

            c.Set("user_id", claims.Subject)
            c.Set("role", claims.Role)
            c.Set("access_jti", claims.ID)
            c.Set("access_exp", claims.ExpiresAt.Time)

            return next(c)
        }
    }
}
