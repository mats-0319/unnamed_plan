package middleware

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

// token structure: `hex({"user_name":[xxx]...}).hex(hash(payload, key))`
// more: doc/design.md 接口访问令牌

type AccessToken struct {
	UserName   string `json:"user_name"`
	ExpireTime int64  `json:"expire_time"`
}

func GenAccessToken(userName string) string {
	tokenBytes, err := json.Marshal(&AccessToken{
		UserName:   userName,
		ExpireTime: time.Now().Add(time.Hour * 6).UnixMilli(), // hard code 'expire time' = 6h
	})
	if err != nil {
		e := ErrServerInternalError().WithCause(err)
		mlog.Error(e.String())
		return ""
	}

	tokenHex := hex.EncodeToString(tokenBytes)

	return fmt.Sprintf("%s.%s", tokenHex, genTokenHash(tokenHex))
}

// OptionalVerifyAccessToken 用于是否登录均可的接口，访问token验证失败不视为错误。
// 如果验证成功，则后续处理函数能拿到'ctx.UserName'
func OptionalVerifyAccessToken(ctx *mhttp.Context) *Error {
	_ = verifyAccessToken(ctx)
	return nil
}

func VerifyAccessToken(ctx *mhttp.Context) (e *Error) {
	e = verifyAccessToken(ctx)
	if e != nil {
		mlog.Error(e.String())
	}

	return
}

func verifyAccessToken(ctx *mhttp.Context) (e *Error) {
	tokenSplit := strings.Split(ctx.AccessToken, ".")
	if len(tokenSplit) != 2 {
		e = ErrInvalidAccessToken().WithParam("token", ctx.AccessToken)
		return
	}

	if hash := genTokenHash(tokenSplit[0]); hash != tokenSplit[1] {
		e = ErrWrongToken().WithParam("token", ctx.AccessToken)
		return
	}

	tokenBytes, err := hex.DecodeString(tokenSplit[0])
	if err != nil {
		e = ErrInvalidAccessToken().WithCause(err)
		return
	}

	token := &AccessToken{}
	if err := json.Unmarshal(tokenBytes, token); err != nil {
		e = ErrInvalidAccessToken().WithCause(err)
		return
	}

	if token.ExpireTime < time.Now().UnixMilli() {
		e = ErrTokenExpired().WithParam("expire time", token.ExpireTime)
		return
	}

	ctx.UserName = token.UserName

	return
}

var hmacKey = GenerateRandomBytes[[]byte](10)

func genTokenHash(tokenHex string) string {
	return HmacSHA256(tokenHex, hmacKey)
}
