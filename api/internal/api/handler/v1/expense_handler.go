package v1

import (
	"github.com/Kudryavkaz/sztuea-api/internal/context/expense"
	"github.com/gofiber/fiber/v2"
)

const (
	expensePrefix = "/expense"
)

func InitExpenseRouter(router fiber.Router) {
	expenseRouterGroup := router.Group(expensePrefix)
	expenseRouterGroup.Use(Authorization)
	expenseRouterGroup.Get("/data", pullExpense)
	expenseRouterGroup.Get("/table", queryTable)
	expenseRouterGroup.Get("/timeline", queryTimeLine)
	expenseRouterGroup.Get("/metric", queryMetric)
}

func pullExpense(ctx *fiber.Ctx) (err error) {
	context := expense.NewContext(ctx, 0)

	context.AddDeferHandler(context.SendResponse)

	context.AddBaseHandler(context.ParsePullRequest).AddBaseHandler(context.CheckPullFields)

	context.AddBaseHandler(context.Pull)

	context.Run()

	return
}

func queryTable(ctx *fiber.Ctx) (err error) {
	context := expense.NewContext(ctx, 0)

	context.AddDeferHandler(context.SendResponse)

	context.AddBaseHandler(context.ParseQueryRequest).AddBaseHandler(context.CheckQueryFields)

	context.AddBaseHandler(context.QueryTable)

	context.Run()

	return
}

func queryTimeLine(ctx *fiber.Ctx) (err error) {
	context := expense.NewContext(ctx, 0)

	context.AddDeferHandler(context.SendResponse)

	context.AddBaseHandler(context.ParseQueryRequest).AddBaseHandler(context.CheckQueryFields)

	context.AddBaseHandler(context.QueryTimeLine)

	context.Run()

	return
}

func queryMetric(ctx *fiber.Ctx) (err error) {
	context := expense.NewContext(ctx, 0)

	context.AddDeferHandler(context.SendResponse)

	context.AddBaseHandler(context.ParseQueryRequest).AddBaseHandler(context.CheckQueryFields)

	context.Run()

	return
}
