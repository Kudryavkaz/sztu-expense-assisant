package grpcclient

import (
	"context"
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/config"
	"github.com/Kudryavkaz/sztuea-api/internal/log"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type Claims struct {
	UserID uint  `json:"user_id"`
	Exp    int64 `json:"exp"`
	*jwt.RegisteredClaims
}

func generateJwtToken(userID uint) (jwtSecret string, err error) {
	claims := Claims{
		UserID: userID,
		Exp:    time.Now().Add(time.Duration(1 * time.Minute)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Logger().Info("jwtSecret", zap.String("jwtSecret", config.Config.GetString("crawler.jwt.secret.key")))
	return token.SignedString([]byte(config.Config.GetString("crawler.jwt.secret.key")))
}

func generateTokenMD(jwtSecret string) metadata.MD {
	return metadata.Pairs("authorization", jwtSecret)
}

func GenerateGrpcCtx(ctx context.Context, userID uint) (grpcCtx context.Context, err error) {
	jwtSecret, err := generateJwtToken(userID)
	if err != nil {
		return
	}
	authorizationMD := generateTokenMD(jwtSecret)

	MD := metadata.Join(authorizationMD)
	grpcCtx = metadata.NewOutgoingContext(ctx, MD)
	return
}
