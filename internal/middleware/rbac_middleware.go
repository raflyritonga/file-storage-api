package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RequireRoles(roles ...string) echo.MiddlewareFunc {
    accepted := map[string]struct{}{}
    for _, role := range roles {
        accepted[role] = struct{}{}
    }

    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            role, _ := c.Get("role").(string)
            if _, ok := accepted[role]; !ok {
                return echo.NewHTTPError(http.StatusForbidden, "forbidden")
            }
            return next(c)
        }
    }
}
