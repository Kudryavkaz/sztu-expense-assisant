package api

import (
	"fmt"
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/context/api"
	"github.com/Kudryavkaz/sztuea-api/internal/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
)

func StartSever(port uint16, prefork bool) (err error) {
	app := fiber.New(fiber.Config{
		AppName:           "github.com/Kudryavkaz/sztuea-api",
		CaseSensitive:     true,
		EnablePrintRoutes: true,
		Prefork:           prefork,
		ErrorHandler:      errHandler,
		IdleTimeout:       15 * time.Minute,
		ReadBufferSize:    8192,
	})

	app.Use(logger.New())

	app.Use(logger.New(logger.Config{
		Format:     "${time}    ${method}${path} ${status} ${respHeader:X-Request-ID} - client: ${ip}:${port} (${ua}) - latency: ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	app.Use(requestid.New())

	app.Use(recover.New())

	v1 := app.Group("/v1")
	InitRouter(v1)

	err = app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Logger().Error("Failed to start server.", zap.String("error", err.Error()))
	}

	return
}

func errHandler(ctx *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	code := -1

	switch e := err.(type) {
	case *fiber.Error:
		status = e.Code
	case *api.Error:
		status = e.Status
		code = e.Code
	}

	// 发送自定义错误
	err = ctx.Status(status).JSON(fiber.Map{
		"code": code,
		"msg":  err.Error(),
	})

	return nil
}

func InitRouter(router fiber.Router) {
}
