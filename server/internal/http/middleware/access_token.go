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

// token structure: `hex({"user":[xxx]...}).hex(hash(payload, key))`
// more: doc/design.md 接口访问令牌

type AccessToken struct {
	User       string `json:"user"`
	ExpireTime int64  `json:"expire_time"`
}

func GenAccessToken(userName string) string {
	tokenBytes, err := json.Marshal(&AccessToken{
		User:       userName,
		ExpireTime: time.Now().Add(time.Hour * 6).UnixMilli(), // hard code 'expire time' = 6h
	})
	if err != nil {
		e := ErrServerInternalError().WithCause(err)
		mlog.Log(e.String())
		return ""
	}

	tokenHex := hex.EncodeToString(tokenBytes)

	return fmt.Sprintf("%s.%s", tokenHex, genTokenHash(tokenHex))
}

func VerifyAccessToken(ctx *mhttp.Context) *Error {
	tokenSplit := strings.Split(ctx.AccessToken, ".")
	if len(tokenSplit) != 2 {
		e := ErrInvalidAccessToken().WithParam("token", ctx.AccessToken)
		mlog.Log(e.String())
		return e
	}

	hash := genTokenHash(tokenSplit[0])
	if hash != tokenSplit[1] {
		e := ErrTokenTamperedWith().WithParam("token", ctx.AccessToken)
		mlog.Log(e.String())
		return e
	}

	tokenBytes, err := hex.DecodeString(tokenSplit[0])
	if err != nil {
		e := ErrInvalidAccessToken().WithCause(err)
		mlog.Log(e.String())
		return e
	}

	token := &AccessToken{}
	err = json.Unmarshal(tokenBytes, token)
	if err != nil {
		e := ErrInvalidAccessToken().WithCause(err)
		mlog.Log(e.String())
		return e
	}

	now := time.Now().UnixMilli()
	if token.ExpireTime < now {
		e := ErrTokenExpired().WithParam("expire time", token.ExpireTime).WithParam("now", now)
		mlog.Log(e.String())
		return e
	}

	ctx.User = token.User

	return nil
}

var hmacKey = GenerateRandomBytes[[]byte](10)

func genTokenHash(tokenHex string) string {
	return HmacSHA256(tokenHex, hmacKey)
}
