package v1

import (
	"github.com/Kudryavkaz/sztuea-api/internal/context/user"
	"github.com/gofiber/fiber/v2"
)

const (
	UserPrefix = "/user"
)

func InitUserRouter(router fiber.Router) {
	userRouterGroup := router.Group(UserPrefix)
	userRouterGroup.Use(Authorization)
	userRouterGroup.Post("", setUserInfo)
}

func setUserInfo(ctx *fiber.Ctx) (err error) {
	context := user.NewContext(ctx, 0)

	context.AddDeferHandler(context.SendResponse)

	context.AddBaseHandler(context.ParseRequest).AddBaseHandler(context.CheckSetUserInfoFields)

	context.AddBaseHandler(context.SetUserInfo)

	context.Run()

	return
}
