package handlers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	"github.com/mats0319/unnamed_plan/server/internal/http/middleware"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/password"
)

func Login(ctx *mhttp.Context) {
	req := &api.LoginReq{}
	if !ctx.ParseParams(req) {
		return
	}

	if len(req.UserName) < 1 || len(req.Password) < 1 {
		e := utils.ErrInvalidParams().WithParam("user name", req.UserName).WithParam("password", req.Password)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	user, e := dal.GetUser(req.UserName)
	if e != nil {
		ctx.ResData = e
		return
	}

	if e := password.VerifyPassword(req.Password, user.Password); e != nil {
		ctx.ResData = e
		return
	}

	//if user.Enable2FA {
	//	if e := verifyTotpCode(req.TotpCode, user.TotpKey); e != nil {
	//		ctx.ResData = e
	//		return
	//	}
	//}

	if e := dal.UpdateUser(user); e != nil { // modify user.UpdatedAt
		ctx.ResData = e
		return
	}

	ctx.Writer.Header().Set(utils.HttpHeader_AccessToken, middleware.GenAccessToken(user.UserName))

	ctx.ResData = &api.LoginRes{
		UserName:  user.UserName,
		Nickname:  user.Nickname,
		IsAdmin:   user.IsAdmin,
		Enable2FA: user.Enable2FA,
	}
}

// verifyTotpCode totpKey should be base32 encoded
func verifyTotpCode(code string, totpKey string) *utils.Error {
	if len(code) != 6 {
		e := utils.ErrInvalidTotpCode().WithParam("code", code)
		mlog.Error(e.String())
		return e
	}

	key := make([]byte, 10)
	n, err := base32.StdEncoding.Decode(key, []byte(totpKey))
	if err != nil {
		e := utils.ErrInvalidTotpKey().WithCause(err)
		mlog.Error(e.String())
		return e
	}
	key = key[:n]

	timestep := time.Now().Unix() / 30

	// allow last/current/next totp code
	validCodes := []string{
		calcTotpCode(key, iTob(timestep-1)),
		calcTotpCode(key, iTob(timestep)),
		calcTotpCode(key, iTob(timestep+1)),
	}

	if !slices.Contains(validCodes, code) {
		e := utils.ErrWrongTotpCode().WithParam("code", code)
		mlog.Error(e.String())
		return e
	}

	return nil
}

func calcTotpCode(key []byte, content []byte) string {
	hasher := hmac.New(sha1.New, key)
	hasher.Write(content)
	hmacHash := hasher.Sum(nil)

	offset := int(hmacHash[len(hmacHash)-1] & 0x0f)
	// 算法要求屏蔽最高有效位
	longPassword := int(hmacHash[offset]&0x7f)<<24 |
		int(hmacHash[offset+1])<<16 |
		int(hmacHash[offset+2])<<8 |
		int(hmacHash[offset+3])

	password := longPassword % int(math.Pow10(6))

	return fmt.Sprintf("%06d", password)
}

func iTob(integer int64) []byte {
	byteArr := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		byteArr[i] = byte(integer & 0xff)
		integer = integer >> 8
	}
	return byteArr
}
