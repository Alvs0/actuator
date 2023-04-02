package impl

import (
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const JWTSecret = "mysecret"

type AuthHandler interface {
	Login(ctx echo.Context) error
	Validate(ctx echo.Context) error
}

type authHandler struct {
	accountQuery AccountQuery
}

func NewAuthHandler(accountQuery AccountQuery) AuthHandler {
	return &authHandler{
		accountQuery: accountQuery,
	}
}

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func CreateMiddleware() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(JWTSecret),
	}

	return echojwt.WithConfig(config)
}