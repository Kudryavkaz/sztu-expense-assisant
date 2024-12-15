package user

import (
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/context/api"
	"github.com/gofiber/fiber/v2"
)

type Base struct {
	sztuAccount  string
	sztuPassword string
	userID       uint
}

type Context struct {
	api.Context
	Base
}

func NewContext(ctx *fiber.Ctx, expireDuration time.Duration) (parseCtx Context) {
	parseCtx = Context{
		Context: api.NewContext(ctx, expireDuration),
	}
	return
}
