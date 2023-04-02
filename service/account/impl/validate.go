package impl

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (a *authHandler) Validate(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name

	return ctx.JSON(http.StatusOK, echo.Map{
		"name": name,
	})
}
