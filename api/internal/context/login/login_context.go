package login

import (
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/context/api"
	"github.com/gofiber/fiber/v2"
)

var ()

type Base struct {
	account        string
	password       string
	repeatPassword string
	sztuAccount    string
	sztuPassword   string
	userId         uint
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
