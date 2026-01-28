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

// token structure: `hex({"user_id":[xxx]...}).hex(hash(payload, key))`
// - payload: hex({"user_id":[xxx]...}), hex str before '.'
// - 轮换hmacKey：在token payload添加key版本号，同时维护多个版本的key
// hash algorithm: hmac-sha256

type AccessToken struct {
	UserID     uint  `json:"user_id"`
	ExpireTime int64 `json:"expire_time"`
}

func GenToken(userID uint) string {
	tokenBytes, err := json.Marshal(&AccessToken{
		UserID:     userID,
		ExpireTime: time.Now().Add(time.Hour * 6).UnixMilli(), // hard code 'expire time' = 6h
	})
	if err != nil {
		e := NewError(ET_ServerInternalError, ED_JsonMarshal).WithCause(err)
		mlog.Log(e.String())
		return ""
	}

	tokenHex := hex.EncodeToString(tokenBytes)

	return fmt.Sprintf("%s.%s", tokenHex, genTokenHash(tokenHex))
}

func VerifyToken(ctx *mhttp.Context) *Error {
	tokenSplit := strings.Split(ctx.AccessToken, ".")
	if len(tokenSplit) != 2 {
		e := NewError(ET_UnauthorizedError, ED_InvalidAccessToken).WithParam("token", ctx.AccessToken)
		mlog.Log(e.String())
		return e
	}

	tokenBytes, err := hex.DecodeString(tokenSplit[0])
	if err != nil {
		e := NewError(ET_UnauthorizedError, ED_InvalidAccessToken).WithCause(err)
		mlog.Log(e.String())
		return e
	}

	hash := genTokenHash(tokenSplit[0])
	if hash != tokenSplit[1] {
		e := NewError(ET_UnauthorizedError, ED_TokenTamperedWith).WithParam("token", ctx.AccessToken)
		mlog.Log(e.String())
		return e
	}

	token := &AccessToken{}
	err = json.Unmarshal(tokenBytes, token)
	if err != nil {
		e := NewError(ET_UnauthorizedError, ED_InvalidAccessToken).WithCause(err)
		mlog.Log(e.String())
		return e
	}

	if token.ExpireTime < time.Now().UnixMilli() {
		e := NewError(ET_UnauthorizedError, ED_TokenExpired).WithParam("expire time", token.ExpireTime)
		mlog.Log(e.String())
		return e
	}

	ctx.UserID = token.UserID

	return nil
}

var hmacKey = GenerateRandomBytes[[]byte](10)

func genTokenHash(tokenHex string) string {
	return HmacSHA256(tokenHex, hmacKey)
}
