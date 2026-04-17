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

// 关于接口应如何使用验证中间件的说明：
// 1. 必须登录后可调用的接口：使用验证中间件，例如写小纸条
// 2. 是否登录均可调用，对访客/用户的处理逻辑相同的接口：不使用中间件，例如注册、登录
// 3. 是否登录均可调用，对访客/用户的处理逻辑不同的接口：使用可选验证中间件，例如上传游戏成绩

type Token struct {
	UserName   string    `json:"user_name"`
	Type       TokenType `json:"type"`
	ExpireTime int64     `json:"expire_time"`
}

type TokenType int8

const (
	TokenType_ApiAccessToken TokenType = 1
	TokenType_MfaToken                 = 2
)

func GenerateApiAccessToken(userName string) string {
	return generateToken(&Token{
		UserName:   userName,
		Type:       TokenType_ApiAccessToken,
		ExpireTime: time.Now().Add(time.Hour * 12).UnixMilli(), // hard code 'expire time' = 12h
	})
}

func GenerateMfaToken(userName string) string {
	tokenIns := &Token{
		UserName:   userName,
		Type:       TokenType_MfaToken,
		ExpireTime: time.Now().Add(time.Minute * 5).UnixMilli(), // hard code 'expire time' = 5min
	}

	token := generateToken(tokenIns)

	NewMfaToken(userName, token, tokenIns.ExpireTime)

	return token
}

func generateToken(token *Token) string {
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		e := ErrServerInternalError().WithCause(err)
		mlog.Error(e.String())
		return ""
	}

	tokenHex := hex.EncodeToString(tokenBytes)

	return fmt.Sprintf("%s.%s", tokenHex, genTokenHash(tokenHex))
}

// OptionalVerifyAccessToken 用于访客和用户均可使用的接口，访问token验证失败视为访客、验证成功可以拿到`ctx.userName`。
func OptionalVerifyAccessToken(ctx *mhttp.Context) *Error {
	_ = verifyAccessToken(ctx)
	return nil
}

func VerifyAccessToken(ctx *mhttp.Context) (e *Error) {
	if e = verifyAccessToken(ctx); e != nil {
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
		e = ErrWrongAccessTokenHash().WithParam("token", ctx.AccessToken)
		return
	}

	tokenBytes, err := hex.DecodeString(tokenSplit[0])
	if err != nil {
		e = ErrDecodeAccessToken().WithCause(err)
		return
	}

	token := &Token{}
	if err := json.Unmarshal(tokenBytes, token); err != nil {
		e = ErrDeserializeAccessToken().WithCause(err)
		return
	}

	if token.Type != TokenType_ApiAccessToken {
		e = ErrInvalidAccessToken().WithParam("token type", token.Type)
		return
	}

	if token.ExpireTime < time.Now().UnixMilli() {
		e = ErrAccessTokenExpired().WithParam("expire time", token.ExpireTime)
		return
	}

	ctx.UserName = token.UserName

	return
}

var hmacKey = GenerateRandomBytes[[]byte](10)

func genTokenHash(tokenHex string) string {
	return HmacSHA256(tokenHex, hmacKey)
}
