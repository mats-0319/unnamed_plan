package middleware

import (
	"errors"
	"time"

	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
)

type accessToken struct {
	Token      string
	ExpireTime int64
}

var accessTokens = map[uint]*accessToken{} // user id - token

func SetToken(id uint, token string) {
	accessTokens[id] = &accessToken{
		Token:      token,
		ExpireTime: time.Now().Add(time.Hour * 6).UnixMilli(), // hard code 'expire time' = 6h
	}
}

func VerifyToken(ctx *mhttp.Context) error {
	t, ok := accessTokens[ctx.UserID]
	if !ok || t.Token != ctx.AccessToken {
		return errors.New("invalid user id or token")
	}

	if t.ExpireTime < time.Now().UnixMilli() {
		return errors.New("token expired")
	}

	return nil
}
