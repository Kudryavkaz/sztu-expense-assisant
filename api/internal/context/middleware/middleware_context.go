package middleware

import (
	"errors"
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/auth"
	"github.com/Kudryavkaz/sztuea-api/internal/context/api"
	"github.com/Kudryavkaz/sztuea-api/internal/resource/database/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Context struct {
	api.Context
}

func NewContext(ctx *fiber.Ctx, expireDuration time.Duration) (parseCtx Context) {
	parseCtx = Context{
		Context: api.NewContext(ctx, expireDuration),
	}
	return
}

func AuthJwt(tokenString string) (userID uint, err error) {
	userID, err = auth.ValidateToken(tokenString)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return userID, errors.New("Login expired, please login again.")
		}
		return userID, err
	}

	if _, err = model.GetUserByID(userID); err != nil {
		return userID, errors.New("User not found.")
	}

	return userID, nil
}
