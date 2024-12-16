package expense

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/context/api"
	"github.com/Kudryavkaz/sztuea-api/internal/grpcclient"
	"github.com/Kudryavkaz/sztuea-api/internal/log"
	"github.com/Kudryavkaz/sztuea-api/internal/resource/cache"
	"github.com/Kudryavkaz/sztuea-api/internal/resource/database/model"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Table struct {
	ExpenseInfoList []SimpleExpense `json:"expense_info_list"`
	Total           int64           `json:"total"`
}

type SimpleExpense struct {
	Action      string  `json:"action"`
	ExpenseTime string  `json:"expense_time"`
	Place       string  `json:"place"`
	Amount      float64 `json:"amount"`
}

type TimeLineChart struct {
	ExpenseInfoList []SimpleTimeExpenseInfo `json:"expense_info_list"`
}

type SimpleTimeExpenseInfo struct {
	ExpenseTime string  `json:"expense_time"`
	Amount      float64 `json:"amount"`
}

type Base struct {
	startTime    int64
	endTime      int64
	place        string
	action       string
	page         int
	perPage      int
	sztuAccount  string
	sztuPassword string
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

type ExpenseResponse struct {
	Ret bool   `json:"ret"`
	Msg string `json:"msg"`
	Obj Object `json:"obj"`
}

type Object struct {
	Expenditure float64       `json:"expenditure"`
	Income      int           `json:"income"`
	Total       int           `json:"total"`
	TotalPage   int           `json:"totalPage"`
	List        []ExpenseInfo `json:"list"`
	PageParam   Page          `json:"pageParam"`
}

type ExpenseInfo struct {
	Sno            string  `json:"sno"`
	TranTypeDesc   string  `json:"tranTypeDesc"`
	TranMethodDesc string  `json:"tranMethodDesc"`
	OriBalance     float64 `json:"oriBalance"`
	Amount         float64 `json:"amount"`
	Balance        float64 `json:"balance"`
	ItemName       string  `json:"itemName"`
	FinishTime     int64   `json:"finishTime"`
	StrFinishTime  string  `json:"strFinishTime"`
	StatusStr      string  `json:"statusStr"`
}

type Page struct {
	PageNum    int `json:"pageNum"`
	NumPerPage int `json:"numPerPage"`
}

type ExpenseRequest struct {
	from             string
	token            string
	yearMonth        string
	pageNum          int
	numPerPage       int
	queryTotalAmount bool
}

const expenseUrlTemplate = "https://sjdxykt.sztu.edu.cn/user/user/account/getAccountRecordList?from=%s&token=%s&yearMonth=%s&pageNum=%d&numPerPage=%d&queryTotalAmount=%t"

func SendRequest(reqest ExpenseRequest) (response ExpenseResponse, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	url := fmt.Sprintf(expenseUrlTemplate, reqest.from, reqest.token, reqest.yearMonth, reqest.pageNum, reqest.numPerPage, reqest.queryTotalAmount)
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{}
	if err = client.Do(req, resp); err != nil {
		log.Logger().Error("[SendRequest]", zap.Error(err))
		return
	}

	statusCode := resp.StatusCode()
	if statusCode == fasthttp.StatusOK {
		log.Logger().Info("[SendRequest] Success", zap.Int("statusCode", statusCode))
	} else {
		log.Logger().Error("[SendRequest] Fail", zap.Int("statusCode", statusCode))
	}

	if err = json.Unmarshal(resp.Body(), &response); err != nil {
		log.Logger().Error("[SendRequest] Fail", zap.Error(err))
		return
	}

	return
}

func GetCookie(account string, password string) (cookie string, err error) {
	cookie, err = cache.Rdb.Get(context.Background(), fmt.Sprintf("cookie:%s", account)).Result()
	if err != nil && err != redis.Nil {
		return
	} else if err == redis.Nil {
		cookie, err = GetCookieFromGrpc(account, password)
	} else {
		if _, err = SendRequest(ExpenseRequest{
			token:            cookie,
			yearMonth:        "2024-12",
			pageNum:          1,
			numPerPage:       1,
			queryTotalAmount: false,
		}); err != nil {
			cookie, err = GetCookieFromGrpc(account, password)
			if err != nil {
				return
			}
		}
	}

	return
}

func GetCookieFromGrpc(sztuAccount string, sztuPassword string) (cookie string, err error) {
	cookie, err = grpcclient.GetCookie(1, sztuAccount, sztuPassword)

	return
}

func GetCookieFromDB(userID uint) (cookie string, err error) {
	user, err := model.GetUserByID(userID)
	if err != nil {
		return
	}

	if user.Cookie == "" {
		cookie, err = UpdateCookie(userID, user.SztuAccount, user.SztuPassword)
		if err != nil {
			return
		}
	}

	// Check Cookie alive
	if _, err = SendRequest(ExpenseRequest{
		token:            user.Cookie,
		yearMonth:        "2024-12",
		pageNum:          1,
		numPerPage:       1,
		queryTotalAmount: false,
	}); err != nil {
		cookie, err = UpdateCookie(userID, user.SztuAccount, user.SztuPassword)
		if err != nil {
			return
		}
	}
	return
}

func UpdateCookie(userID uint, account string, password string) (cookie string, err error) {
	cookie, err = grpcclient.GetCookie(userID, account, password)
	if err != nil {
		return
	}

	user, err := model.GetUserByID(userID)
	if err != nil {
		return
	}

	user.Cookie = cookie
	err = model.UpdateAccountByID(userID, user)
	if err != nil {
		return
	}
	cache.Rdb.Set(context.Background(), fmt.Sprintf("cookie:%d", userID), cookie, 24*time.Hour)

	return
}

func ToExpenseDO(respExpense ExpenseInfo) (expenseDO *model.Expense) {
	expenseDO = &model.Expense{
		// UserID:         userID,
		Sno:            respExpense.Sno,
		TranTypeDesc:   respExpense.TranTypeDesc,
		TranMethodDesc: respExpense.TranMethodDesc,
		OriBalance:     respExpense.OriBalance,
		Amount:         respExpense.Amount,
		Balance:        respExpense.Balance,
		ItemName:       respExpense.ItemName,
		FinishTime:     respExpense.FinishTime,
		StrFinishTime:  respExpense.StrFinishTime,
		StatusStr:      respExpense.StatusStr,
	}
	return
}
