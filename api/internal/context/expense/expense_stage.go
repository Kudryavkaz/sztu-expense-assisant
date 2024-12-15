package expense

import (
	"context"
	"fmt"
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/context/api"
	"github.com/Kudryavkaz/sztuea-api/internal/lock"
	"github.com/Kudryavkaz/sztuea-api/internal/log"
	"github.com/Kudryavkaz/sztuea-api/internal/resource/database/model"
	"go.uber.org/zap"
)

func (expenseCtx *Context) ParseRequest(ctx context.Context) (needBreak bool, errMsg string, err error) {
	body := expenseCtx.FCtx.Body()
	if body != nil {
		needBreak = true
		expenseCtx.APIError = api.ErrParseRequest
		return
	}

	expenseCtx.userID = expenseCtx.FCtx.Locals("userID").(uint)

	return
}

func (expenseCtx *Context) CheckPullFields(ctx context.Context) (needBreak bool, errMsg string, err error) {
	return
}

func (expenseCtx *Context) Pull(ctx context.Context) (needBreak bool, errMsg string, err error) {
	cookie, err := GetCookie(expenseCtx.userID)
	log.Logger().Info("[Pull]", zap.String("cookie", cookie))

	mu, err := lock.GetLock(fmt.Sprintf("userID$%d", expenseCtx.userID), 1, time.Minute)
	if err != nil {
		log.Logger().Error("[Pull] GetLock Fail", zap.Error(err))
		needBreak = true
		expenseCtx.APIError = api.ErrGetLock
		return
	}
	defer mu.Release()

	latestExpenseTime, err := model.GetLatestExpenseByUserID(expenseCtx.userID)
	if err != nil {
		log.Logger().Error("[Pull] GetLatestExpenseByUserID Fail", zap.Error(err))
		needBreak = true
		expenseCtx.APIError = api.ErrQueryExpense
		return
	}

	startMonth := time.Unix(latestExpenseTime, 0).Format("2006-01")
	endMonth := time.Now().Format("2006-01")

	startDate, _ := time.Parse("2006-01", startMonth)
	endDate, _ := time.Parse("2006-01", endMonth)

	var resp ExpenseResponse
	expenses := make(model.Expenses, 0)
	for current := startDate; current.Before(endDate) || current.Equal(endDate); current = current.AddDate(0, 1, 0) {
		log.Logger().Info("[Pull]", zap.String("current", current.Format("2006-01")))

		for pageNum := 1; ; pageNum++ {
			resp, err = SendRequest(ExpenseRequest{
				token:      cookie,
				yearMonth:  current.Format("2006-01"),
				pageNum:    pageNum,
				numPerPage: 100,
			})
			for _, expense := range resp.Obj.List {
				if expense.FinishTime > latestExpenseTime {
					expenses = append(expenses, ToExpenseDO(expense, expenseCtx.userID))
				}
			}
			if pageNum == resp.Obj.TotalPage {
				break
			}
		}
	}
	if err = expenses.Create(); err != nil {
		log.Logger().Error("[Pull] Create Fail", zap.Error(err))
		needBreak = true
		expenseCtx.APIError = api.ErrCreateExpense
		return
	}

	return
}
