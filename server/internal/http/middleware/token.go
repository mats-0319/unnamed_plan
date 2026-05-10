package middleware

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

type Token struct {
	UserName   string    `json:"user_name"`
	Type       TokenType `json:"type"`
	ExpireTime int64     `json:"expire_time"`
}

type TokenType int8

const (
	TokenType_APIAccessToken TokenType = 1
	TokenType_MFAToken                 = 2
)

func GenerateAPIAccessToken(userName string) string {
	return generateToken(&Token{
		UserName:   userName,
		Type:       TokenType_APIAccessToken,
		ExpireTime: time.Now().Add(time.Hour * 12).UnixMilli(), // hard code 'expire time' = 12h
	})
}

func GenerateMFAToken(userName string) string {
	tokenIns := &Token{
		UserName:   userName,
		Type:       TokenType_MFAToken,
		ExpireTime: time.Now().Add(time.Minute * 5).UnixMilli(), // hard code 'expire time' = 5min
	}

	token := generateToken(tokenIns)

	NewMFAToken(userName, token, tokenIns.ExpireTime)

	return token
}

func generateToken(token *Token) string {
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		e := utils.ErrServerInternalError().WithCause(err)
		mlog.Error(e.String())
		return ""
	}

	tokenHex := hex.EncodeToString(tokenBytes)

	return fmt.Sprintf("%s.%s", tokenHex, genTokenHash(tokenHex))
}

// OptionalVerifyAccessToken 用于访客和用户均可使用的接口，访问token验证失败视为访客、验证成功可以拿到`ctx.userName`。
func OptionalVerifyAccessToken(ctx *mhttp.Context) *utils.Error {
	_ = verifyAccessToken(ctx)
	return nil
}

func VerifyAccessToken(ctx *mhttp.Context) *utils.Error {
	if e := verifyAccessToken(ctx); e != nil {
		mlog.Error(e.String())
		return e
	}

	return nil
}

// verifyAccessToken 函数内不打印错误，因为部分场景允许验证错误（例如上方的可选验证函数）
func verifyAccessToken(ctx *mhttp.Context) (e *utils.Error) {
	tokenSplit := strings.Split(ctx.AccessToken, ".")
	if len(tokenSplit) != 2 {
		e = utils.ErrInvalidAccessToken().WithParam("token", ctx.AccessToken)
		return
	}

	if hash := genTokenHash(tokenSplit[0]); hash != tokenSplit[1] {
		e = utils.ErrWrongAccessTokenHash().WithParam("token", ctx.AccessToken)
		return
	}

	tokenBytes, err := hex.DecodeString(tokenSplit[0])
	if err != nil {
		e = utils.ErrDecodeAccessToken().WithCause(err)
		return
	}

	token := &Token{}
	if err := json.Unmarshal(tokenBytes, token); err != nil {
		e = utils.ErrDeserializeAccessToken().WithCause(err)
		return
	}

	if token.Type != TokenType_APIAccessToken {
		e = utils.ErrInvalidAccessToken().WithParam("token type", token.Type)
		return
	}

	if token.ExpireTime < time.Now().UnixMilli() {
		e = utils.ErrAccessTokenExpired().WithParam("expire time", token.ExpireTime)
		return
	}

	ctx.UserName = token.UserName

	return
}

var hmacKey = utils.GenerateRandomBytes[[]byte](10)

func genTokenHash(tokenHex string) string {
	return utils.HmacSHA256(tokenHex, hmacKey)
}
