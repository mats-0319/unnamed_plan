package middleware

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

// token structure: `{"user_id":[xxx],"expire_time":[xxx]}.[hash]`
// hash: sha256([random str]+`{[xxx]}`)

type AccessToken struct {
	UserID     uint  `json:"user_id"`
	ExpireTime int64 `json:"expire_time"`
}

var hashSalt = GenerateRandomStr(10)

func GenToken(userID uint) string {
	jsonBytes, err := json.Marshal(&AccessToken{
		UserID:     userID,
		ExpireTime: time.Now().Add(time.Hour * 6).UnixMilli(), // hard code 'expire time' = 6h
	})
	if err != nil {
		e := NewError(ET_ServerInternalError, ED_JsonMarshal).WithCause(err)
		mlog.Log(e.String())
		return ""
	}

	return fmt.Sprintf("%s.%s", jsonBytes, genTokenHash(string(jsonBytes)))
}

func VerifyToken(ctx *mhttp.Context) *Error {
	index := strings.LastIndex(ctx.AccessToken, ".")
	tokenSplit := []string{ctx.AccessToken[:index], ctx.AccessToken[index+1:]}

	hash := genTokenHash(tokenSplit[0])
	if hash != tokenSplit[1] {
		e := NewError(ET_UnauthorizedError, ED_TokenTamperedWith).WithParam("token", ctx.AccessToken)
		mlog.Log(e.String())
		return e
	}

	token := &AccessToken{}
	err := json.Unmarshal([]byte(tokenSplit[0]), token)
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

func genTokenHash(tokenStr string) string {
	return CalcSHA256(hashSalt, tokenStr)
}
