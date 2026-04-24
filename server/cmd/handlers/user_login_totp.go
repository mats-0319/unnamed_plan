package handlers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"slices"
	"time"

	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	"github.com/mats0319/unnamed_plan/server/internal/http/middleware"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func LoginMFA(ctx *mhttp.Context) {
	req := &api.LoginMFAReq{}
	if !ctx.ParseParams(req) {
		return
	}

	if len(req.MFAToken) < 1 || len(req.TOTPCode) < 1 {
		e := utils.ErrInvalidParams().WithParam("mfa token", req.MFAToken).WithParam("totp code", req.TOTPCode)
		ctx.ResData = e
		mlog.Error(e.String())
		return
	}

	userName, e := middleware.DecodeUserNameFromMFAToken(req.MFAToken)
	if e != nil {
		ctx.ResData = e
		return
	}

	if e := middleware.VerifyMFAToken(userName, req.MFAToken); e != nil {
		ctx.ResData = e
		return
	}

	user, e := dal.GetUser(userName)
	if e != nil {
		ctx.ResData = e
		return
	}

	if e := verifyTOTPCode(req.TOTPCode, user.TOTPKey); e != nil {
		ctx.ResData = e
		return
	}

	_ = dal.UpdateUser(user) // modify user.UpdatedAt

	ctx.Writer.Header().Set(utils.HTTPHeader_AccessToken, middleware.GenerateAPIAccessToken(user.UserName))

	ctx.ResData = &api.LoginMFARes{
		UserName:  user.UserName,
		Nickname:  user.Nickname,
		IsAdmin:   user.IsAdmin,
		EnableMFA: user.EnableMFA,
	}
}

// verifyTOTPCode totpKey should be base32 encoded
func verifyTOTPCode(code string, totpKey string) *utils.Error {
	if len(code) != 6 {
		e := utils.ErrInvalidTOTPCode().WithParam("code", code)
		mlog.Error(e.String())
		return e
	}

	key := make([]byte, 10)
	n, err := base32.StdEncoding.Decode(key, []byte(totpKey))
	if err != nil {
		e := utils.ErrInvalidTOTPKey().WithCause(err)
		mlog.Error(e.String())
		return e
	}
	key = key[:n]

	timestep := time.Now().Unix() / 30

	// allow last/current/next totp code
	validCodes := []string{
		calcTOTPCode(key, iTob(timestep-1)),
		calcTOTPCode(key, iTob(timestep)),
		calcTOTPCode(key, iTob(timestep+1)),
	}

	if !slices.Contains(validCodes, code) {
		e := utils.ErrWrongTOTPCode().WithParam("code", code)
		mlog.Error(e.String())
		return e
	}

	return nil
}

func calcTOTPCode(key []byte, content []byte) string {
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
	byteSlice := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		byteSlice[i] = byte(integer & 0xff)
		integer = integer >> 8
	}
	return byteSlice
}
