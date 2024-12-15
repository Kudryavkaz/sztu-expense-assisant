package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/Kudryavkaz/sztuea-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Authorized bool `json:"authorized"`
	UserID     uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(id uint) (string, error) {
	claims := Claims{
		Authorized: true,
		UserID:     id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Config.GetInt("middleware.token_expired")) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.GetString("middleware.jwt.secret")))
}

func ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.Config.GetString("middleware.jwt.secret")), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok && !token.Valid {
		return 0, errors.New("token is invalid")
	}

	return claims.UserID, nil
}

func JudgeJwt(authHeader string) (ok bool) {
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) == 2 && authHeaderParts[0] == "Bearer" {
		return true
	}
	return false
}
