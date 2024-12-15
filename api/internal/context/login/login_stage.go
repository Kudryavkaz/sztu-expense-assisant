package login

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"

	"github.com/Kudryavkaz/sztuea-api/internal/auth"
	"github.com/Kudryavkaz/sztuea-api/internal/context/api"
	"github.com/Kudryavkaz/sztuea-api/internal/resource/database/model"
)

func (loginCtx *Context) ParseRequest(ctx context.Context) (needBreak bool, errMsg string, err error) {
	var request struct {
		Account        string `json:"account"`
		Password       string `json:"password"`
		RepeatPassword string `json:"repeat_password"`
	}

	body := loginCtx.FCtx.Body()
	if err = json.Unmarshal(body, &request); err != nil {
		needBreak = true
		loginCtx.APIError = api.ErrParseRequest
		return
	}

	loginCtx.account = request.Account
	loginCtx.password = request.Password
	loginCtx.repeatPassword = request.RepeatPassword

	return
}

func (userCtx *Context) CheckLoginField(ctx context.Context) (needBreak bool, errMsg string, err error) {
	if userCtx.account == "" || userCtx.password == "" {
		needBreak = true
		userCtx.APIError = api.ErrMissToken
		return
	}

	return
}

func (userCtx *Context) CheckRegisterField(ctx context.Context) (needBreak bool, errMsg string, err error) {
	if userCtx.account == "" || userCtx.password == "" {
		needBreak = true
		userCtx.APIError = api.ErrMissToken
		return
	}

	if userCtx.password != userCtx.repeatPassword {
		needBreak = true
		userCtx.APIError = api.ErrRepeatPassword
		return
	}

	return
}

func (userCtx *Context) Login(ctx context.Context) (needBreak bool, errMsg string, err error) {
	h := md5.New()
	io.WriteString(h, userCtx.password)

	passwordMd5 := string(fmt.Sprintf("%x", h.Sum(nil)))

	user, err := model.GetUserByAccount(userCtx.account)
	if err != nil {
		needBreak = true
		userCtx.APIError = api.ErrQueryUser
		return
	}

	if user.Password != passwordMd5 {
		needBreak = true
		userCtx.APIError = api.ErrPassword
		return
	}
	userCtx.userId = user.ID

	return
}

func (userCtx *Context) Register(ctx context.Context) (needBreak bool, errMsg string, err error) {
	passwordMd5 := string(fmt.Sprintf("%x", md5.Sum([]byte(userCtx.password))))
	user := &model.User{
		Account:  userCtx.account,
		Password: passwordMd5,
	}

	if _, err = model.GetUserByAccount(user.Account); err == nil {
		needBreak = true
		userCtx.APIError = api.ErrRepeatAccount
		return
	}

	if err = user.Create(); err != nil {
		needBreak = true
		userCtx.APIError = api.ErrCreateUser
		return
	}

	return
}

func (userCtx *Context) GenerateToken(ctx context.Context) (needBreak bool, errMsg string, err error) {
	token, err := auth.GenerateToken(userCtx.userId)
	if err != nil {
		needBreak = true
		userCtx.APIError = api.ErrGenerateToken
		return
	}

	userCtx.Data = map[string]interface{}{
		"token": token,
	}

	return
}
