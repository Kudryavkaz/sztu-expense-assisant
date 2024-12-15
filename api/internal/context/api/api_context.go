package api

import (
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/context"
	"github.com/gofiber/fiber/v2"
)

type Base struct {
	APIError *Error
	Data     any
}

type Context struct {
	context.BaseContext
	FCtx *fiber.Ctx
	Base
}

func NewContext(ctx *fiber.Ctx, expireDuration time.Duration) (apiCtx Context) {
	baseCtx := context.BaseContext{}
	baseCtx.Init(expireDuration)
	apiCtx = Context{
		FCtx:        ctx,
		BaseContext: baseCtx,
		Base:        Base{},
	}
	return
}
