package middleware

import (
	"net/http"

	"github.com/labstack/echo"
)

func Auth(checkAuth bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sc := c.(*SCContext)
			if checkAuth && sc.isAuth() {
				return next(c)
			} else if !checkAuth && sc.isAuth() {
				return &echo.HTTPError{
					Code:    http.StatusBadRequest,
					Message: "You are already logged in",
				}
			} else if !checkAuth && !sc.isAuth() {
				return next(c)
			}
			return &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "You are not logged in",
			}
		}
	}
}
