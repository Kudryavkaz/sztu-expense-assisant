package user

import (
	"context"
	"encoding/json"

	"github.com/Kudryavkaz/sztuea-api/internal/context/api"
	"github.com/Kudryavkaz/sztuea-api/internal/resource/database/model"
)

func (userCtx *Context) ParseRequest(ctx context.Context) (needBreak bool, errMsg string, err error) {
	var request struct {
		SztuAccount  string `json:"sztu_account"`
		SztuPassword string `json:"sztu_password"`
	}

	body := userCtx.FCtx.Body()
	if err = json.Unmarshal(body, &request); err != nil {
		needBreak = true
		userCtx.APIError = api.ErrParseRequest
		return
	}

	userCtx.sztuAccount = request.SztuAccount
	userCtx.sztuPassword = request.SztuPassword
	userCtx.userID = userCtx.FCtx.Locals("userID").(uint)

	return
}

func (userCtx *Context) CheckSetUserInfoFields(ctx context.Context) (needBreak bool, errMsg string, err error) {
	if userCtx.sztuAccount == "" || userCtx.sztuPassword == "" {
		needBreak = true
		userCtx.APIError = api.ErrMissToken
		return
	}

	return
}

func (userCtx *Context) SetUserInfo(ctx context.Context) (needBreak bool, errMsg string, err error) {
	user := model.User{
		SztuAccount:  userCtx.sztuAccount,
		SztuPassword: userCtx.sztuPassword,
	}
	if err = model.UpdateAccountByID(userCtx.userID, user); err != nil {
		needBreak = true
		userCtx.APIError = api.ErrUpdateUser
		return
	}

	return
}
