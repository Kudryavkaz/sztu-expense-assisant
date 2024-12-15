package v1

import (
	"github.com/Kudryavkaz/sztuea-api/internal/context/login"
	"github.com/gofiber/fiber/v2"
)

const (
	LoginPrefix    = "/login"
	RegisterPrefix = "/register"
)

func InitLoginRouter(router fiber.Router) {
	loginRouterGroup := router.Group(LoginPrefix)
	loginRouterGroup.Post("", doLogin)

	registerRouterGroup := router.Group(RegisterPrefix)
	registerRouterGroup.Post("", doRegister)
}

func doLogin(ctx *fiber.Ctx) (err error) {
	context := login.NewContext(ctx, 0)

	context.AddDeferHandler(context.SendResponse)

	context.AddBaseHandler(context.ParseRequest).AddBaseHandler(context.CheckLoginField)

	context.AddBaseHandler(context.Login).AddBaseHandler(context.GenerateToken)

	context.Run()

	return
}

func doRegister(ctx *fiber.Ctx) (err error) {
	context := login.NewContext(ctx, 0)

	context.AddDeferHandler(context.SendResponse)

	context.AddBaseHandler(context.ParseRequest).AddBaseHandler(context.CheckRegisterField)

	context.AddBaseHandler(context.Register)

	context.Run()

	return
}
