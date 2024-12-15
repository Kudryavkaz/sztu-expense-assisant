package v1

import (
	"github.com/Kudryavkaz/sztuea-api/internal/context/middleware"
	"github.com/gofiber/fiber/v2"
)

func Authorization(ctx *fiber.Ctx) (err error) {
	context := middleware.NewContext(ctx, 0)

	context.AddDeferHandler(context.SendResponse)

	context.AddBaseHandler(context.Auth)

	context.Run()

	return
}
