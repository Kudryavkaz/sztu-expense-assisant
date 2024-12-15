package middleware

import (
	"context"
	"strings"

	"github.com/Kudryavkaz/sztuea-api/internal/auth"
	"github.com/Kudryavkaz/sztuea-api/internal/context/api"
	"github.com/gofiber/fiber/v2"
)

func (middlewareContext *Context) Auth(_ context.Context) (needBreak bool, errMsg string, err error) {
	authHeader := middlewareContext.FCtx.Get("Authorization")
	if authHeader == "" {
		needBreak = true
		middlewareContext.APIError = api.ErrMissToken
		return
	}

	if ok := auth.JudgeJwt(authHeader); !ok {
		needBreak = true
		middlewareContext.APIError = api.ErrToken
		return
	}

	userID, err := AuthJwt(strings.TrimPrefix(authHeader, "Bearer "))
	if err != nil {
		needBreak = true
		middlewareContext.APIError = api.ErrToken
		return
	}

	middlewareContext.FCtx.Locals("userID", userID)
	middlewareContext.FCtx.Next()

	return
}

func (middlewareContext *Context) SendResponse(_ context.Context) (needBreak bool, errMsg string, err error) {
	if middlewareContext.APIError != nil {
		middlewareContext.FCtx.Status(middlewareContext.APIError.Status).JSON(fiber.Map{
			"code": middlewareContext.APIError.Code,
			"msg":  middlewareContext.APIError.Msg,
		})
		return
	}
	return
}
