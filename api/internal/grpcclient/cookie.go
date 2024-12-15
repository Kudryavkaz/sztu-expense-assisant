package grpcclient

import (
	"context"

	"github.com/Kudryavkaz/sztuea-api/internal/grpcclient/protos"
)

func GetCookie(userID uint, account string, password string) (cookie string, err error) {
	ctx, err := GenerateGrpcCtx(context.Background(), userID)
	if err != nil {
		return
	}

	resp, err := Client.GetCookie(ctx, &protos.GetCookieRequest{
		Account:  account,
		Password: password,
	})
	if err != nil {
		return
	}
	cookie = resp.Cookie

	return
}
