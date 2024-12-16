package expense

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/context/api"
	"github.com/Kudryavkaz/sztuea-api/internal/lock"
	"github.com/Kudryavkaz/sztuea-api/internal/log"
	"github.com/Kudryavkaz/sztuea-api/internal/resource/database/model"
	"go.uber.org/zap"
)

func (expenseCtx *Context) ParsePullRequest(ctx context.Context) (needBreak bool, errMsg string, err error) {
	var request struct {
		SztuAccount  string `json:"sztu_account"`
		SztuPassword string `json:"sztu_password"`
	}
	body := expenseCtx.FCtx.Body()
	if err = json.Unmarshal(body, &request); err != nil {
		log.Logger().Error("[ParseQueryRequest] Unmarshal Fail", zap.Error(err))
		needBreak = true
		expenseCtx.APIError = api.ErrParseRequest
		return
	}

	// expenseCtx.userID = expenseCtx.FCtx.Locals("userID").(uint)
	expenseCtx.sztuAccount = request.SztuAccount
	expenseCtx.sztuPassword = request.SztuPassword

	return
}

func (expenseCtx *Context) CheckPullFields(ctx context.Context) (needBreak bool, errMsg string, err error) {
	if expenseCtx.sztuAccount == "" || expenseCtx.sztuPassword == "" {
		needBreak = true
		expenseCtx.APIError = api.ErrParseRequest
		return
	}

	return
}

func (expenseCtx *Context) Pull(ctx context.Context) (needBreak bool, errMsg string, err error) {
	cookie, err := GetCookie(expenseCtx.sztuAccount, expenseCtx.sztuPassword)
	log.Logger().Info("[Pull]", zap.String("cookie", cookie))

	mu, err := lock.GetLock(fmt.Sprintf("sno$%s", expenseCtx.sztuAccount), 1, time.Minute)
	if err != nil {
		log.Logger().Error("[Pull] GetLock Fail", zap.Error(err))
		needBreak = true
		expenseCtx.APIError = api.ErrGetLock
		return
	}
	defer mu.Release()

	latestExpenseTime, err := model.GetLatestExpenseBySno(expenseCtx.sztuAccount)
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
	retryTime := 5
	for current := startDate; current.Before(endDate) || current.Equal(endDate); current = current.AddDate(0, 1, 0) {
		log.Logger().Info("[Pull]", zap.String("current", current.Format("2006-01")))

		for pageNum := 1; retryTime > 0; pageNum++ {
			resp, err = SendRequest(ExpenseRequest{
				token:      cookie,
				yearMonth:  current.Format("2006-01"),
				pageNum:    pageNum,
				numPerPage: 100,
			})
			if err != nil {
				pageNum--
				retryTime--
			}
			for _, expense := range resp.Obj.List {
				if expense.FinishTime > latestExpenseTime {
					expenses = append(expenses, ToExpenseDO(expense))
				}
			}
			if pageNum == resp.Obj.TotalPage {
				break
			}
		}
	}
	if len(expenses) == 0 {
		return
	}

	if err = expenses.Create(); err != nil {
		log.Logger().Error("[Pull] Create Fail", zap.Error(err))
		needBreak = true
		expenseCtx.APIError = api.ErrCreateExpense
		return
	}

	return
}

func (expenseCtx *Context) ParseQueryRequest(ctx context.Context) (needBreak bool, errMsg string, err error) {
	var request struct {
		SztuAccount string `json:"sztu_account"`
		StartTime   int64  `json:"start_time"`
		EndTime     int64  `json:"end_time"`
		Page        int    `json:"page"`
		PerPage     int    `json:"per_page"`
		Place       string `json:"place"`
		Action      string `json:"action"`
	}

	body := expenseCtx.FCtx.Body()
	if err = json.Unmarshal(body, &request); err != nil {
		log.Logger().Error("[ParseQueryRequest] Unmarshal Fail", zap.Error(err))
		needBreak = true
		expenseCtx.APIError = api.ErrParseRequest
		return
	}

	expenseCtx.sztuAccount = request.SztuAccount
	expenseCtx.startTime = request.StartTime
	expenseCtx.endTime = request.EndTime
	expenseCtx.place = request.Place
	expenseCtx.action = request.Action
	expenseCtx.page = request.Page
	expenseCtx.perPage = request.PerPage
	// expenseCtx.userID = expenseCtx.FCtx.Locals("userID").(uint)

	return
}

func (expenseCtx *Context) CheckQueryFields(ctx context.Context) (needBreak bool, errMsg string, err error) {
	// if expenseCtx.sztuAccount == "" || expenseCtx.startTime == 0 || expenseCtx.endTime == 0 || expenseCtx.startTime > expenseCtx.endTime {
	// 	needBreak = true
	// 	expenseCtx.APIError = api.ErrParseRequest
	// 	return
	// }
	if expenseCtx.sztuAccount == "" {
		needBreak = true
		expenseCtx.APIError = api.ErrParseRequest
		return
	}

	return
}

func (expenseCtx *Context) QueryTable(ctx context.Context) (needBreak bool, errMsg string, err error) {
	expense := &model.Expense{
		Sno:          expenseCtx.sztuAccount,
		ItemName:     expenseCtx.place,
		TranTypeDesc: expenseCtx.action,
	}
	expenses, err := expense.GetExpensesByPage(expenseCtx.page, expenseCtx.perPage)
	if err != nil {
		log.Logger().Error("[QueryTable] GetExpensesByUserIDAndTimeRange Fail", zap.Error(err))
		needBreak = true
		expenseCtx.APIError = api.ErrQueryExpense
		return
	}

	total, err := expense.GetTotalExpense()
	if err != nil {
		log.Logger().Error("[QueryTable] GetTotalExpense Fail", zap.Error(err))
		needBreak = true
		expenseCtx.APIError = api.ErrQueryExpense
		return
	}

	data := Table{
		ExpenseInfoList: make([]SimpleExpense, 0),
		Total:           total,
	}
	for _, expense := range expenses {
		data.ExpenseInfoList = append(data.ExpenseInfoList, SimpleExpense{
			Action:      expense.TranTypeDesc,
			ExpenseTime: time.Unix(expense.FinishTime/1000, 0).Format("2006-01-02 15:04:05"),
			Place:       expense.ItemName,
			Amount:      expense.Amount,
		})
	}

	expenseCtx.Data = data
	return
}

func (expenseCtx *Context) QueryTimeLine(ctx context.Context) (needBreak bool, errMsg string, err error) {
	expense := &model.Expense{
		Sno:          expenseCtx.sztuAccount,
		ItemName:     expenseCtx.place,
		TranTypeDesc: expenseCtx.action,
	}
	timeLine, err := expense.GetExpenseTimeLine(expenseCtx.startTime, expenseCtx.endTime)
	if err != nil {
		log.Logger().Error("[QueryChart] GetExpenseTimeLine Fail", zap.Error(err))
		needBreak = true
		expenseCtx.APIError = api.ErrQueryExpense
		return
	}

	data := TimeLineChart{
		ExpenseInfoList: make([]SimpleTimeExpenseInfo, 0),
	}
	for _, line := range timeLine {
		date, _ := time.Parse(time.RFC3339, line.EventDate)
		data.ExpenseInfoList = append(data.ExpenseInfoList, SimpleTimeExpenseInfo{
			ExpenseTime: date.Format("2006-01-02"),
			Amount:      line.Amount,
		})
	}

	expenseCtx.Data = data
	return
}

func (expenseCtx *Context) QueryMetric(ctx context.Context) (needBreak bool, errMsg string, err error) {

	return
}
