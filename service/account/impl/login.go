package impl

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func (a *authHandler) Login(ctx echo.Context) error {
	username := ctx.FormValue("username")
	password := ctx.FormValue("password")

	userDbs, err := a.accountQuery.GetUser(username)
	if err != nil {
		log.Warn("[Account:Login] user not found")
		return echo.ErrUnauthorized
	}

	if userDbs == nil || len(userDbs) == 0 {
		return echo.ErrUnauthorized
	}

	var found bool
	for _, userDb := range userDbs {
		err = bcrypt.CompareHashAndPassword([]byte(userDb.EncryptedPassword), []byte(password))
		if err != nil {
			continue
		}

		found = true
	}

	if !found {
		return echo.ErrUnauthorized
	}

	claims := &jwtCustomClaims{
		username,
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	authToken, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"authorization_token": authToken,
	})
}
