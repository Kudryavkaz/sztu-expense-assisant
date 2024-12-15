package api

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func (apiCtx *Context) SendResponse(ctx context.Context) (needBreak bool, errMsg string, err error) {
	if apiCtx.APIError != nil {
		apiCtx.FCtx.Status(apiCtx.APIError.Status).JSON(fiber.Map{
			"code": apiCtx.APIError.Code,
			"msg":  apiCtx.APIError.Msg,
		})
		return
	}
	apiCtx.FCtx.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": 0,
		"msg":  "success",
		"data": apiCtx.Data,
	})
	return
}
