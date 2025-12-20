package middleware

import (
	"time"

	. "github.com/mats0319/unnamed_plan/server/internal/const"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
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
		err := NewError(ET_AuthenticateError, ED_InvalidUserIDOrToken).
			WithParam("user id", ctx.UserID).WithParam("token", ctx.AccessToken)
		mlog.Log(err.String())
		return err
	}

	if t.ExpireTime < time.Now().UnixMilli() {
		err := NewError(ET_AuthenticateError, ED_TokenExpired).WithParam("expire time", t.ExpireTime)
		mlog.Log(err.String())
		return err
	}

	return nil
}
