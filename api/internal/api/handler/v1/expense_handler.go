package v1

import "github.com/gofiber/fiber/v2"
import "github.com/Kudryavkaz/sztuea-api/internal/context/expense"

const (
	expensePrefix = "/expense"
)

func InitExpenseRouter(router fiber.Router) {
	expenseRouterGroup := router.Group(expensePrefix)
	expenseRouterGroup.Use(Authorization)
	expenseRouterGroup.Get("", pullExpense)
}

func pullExpense(ctx *fiber.Ctx) (err error) {
	context := expense.NewContext(ctx, 0)

	context.AddDeferHandler(context.SendResponse)

	context.AddBaseHandler(context.ParseRequest).AddBaseHandler(context.CheckPullFields)

	context.AddBaseHandler(context.Pull)

	context.Run()

	return
}
