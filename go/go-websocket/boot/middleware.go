package boot

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// SetupMiddleware Sets ujp the middlewares used
func SetupMiddleware(e *echo.Echo) {

	e.Use(middleware.Logger())
	// e.Use(middleware.RequestID())
	// e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"*"},
	// 	AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	// 	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	// }))
	// e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
	// 	Level: 2,
	// }))
}
