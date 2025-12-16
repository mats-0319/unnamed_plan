package handlers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"time"

	"github.com/mats0319/unnamed_plan/server/1_user/db"
	"github.com/mats0319/unnamed_plan/server/1_user/middleware"
	. "github.com/mats0319/unnamed_plan/server/internal/const"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func Login(ctx *mhttp.Context) {
	req := &api.LoginReq{}
	if !ctx.ParseParams(req) {
		return
	}

	user, err := db.GetUser(req.UserName)
	if err != nil {
		ctx.ResData = err
		return
	}

	if utils.CalcSHA256(req.Password, user.Salt) != user.Password {
		e := NewError(ET_ParamsError, ED_InvalidPwd).WithParam("user name", req.UserName)
		mlog.Log(e.String())
		ctx.ResData = e
		return
	}

	if len(user.TotpKey) > 0 {
		err = verifyTotpCode(req.TotpCode, user.TotpKey)
		if err != nil {
			ctx.ResData = err
			return
		}
	}

	user.LastLogin = time.Now().UnixMilli()
	if err = db.UpdateUser(user); err != nil {
		mlog.Log("update user failed", mlog.Field("error", err))
		ctx.ResData = err
		return
	}

	// token
	token := string(utils.GenerateRandomBytes(10))
	middleware.SetToken(user.ID, token)

	ctx.ResHeaders[HttpHeader_AccessToken] = token

	ctx.ResData = &api.LoginRes{
		ResBase:  api.ResBase{IsSuccess: true},
		UserID:   user.ID,
		UserName: user.UserName,
		Nickname: user.Nickname,
		IsAdmin:  user.IsAdmin,
	}
}

// verifyTotpCode totpKey is base32 encoded
func verifyTotpCode(code string, totpKey string) error {
	if len(code) != 6 {
		e := NewError(ET_ParamsError, ED_InvalidTotpCode).WithParam("code", code)
		mlog.Log(e.String())
		return e
	}

	key := make([]byte, 10)
	n, err := base32.StdEncoding.Decode(key, []byte(totpKey))
	if err != nil {
		e := NewError(ET_ParamsError, ED_Base32Decode).WithCause(err)
		mlog.Log(e.String())
		return e
	}
	key = key[:n]

	timestep := time.Now().Unix() / 30

	validCodes := []string{
		calcTotpCode(key, iTob(timestep-1)),
		calcTotpCode(key, iTob(timestep)),
		calcTotpCode(key, iTob(timestep+1)),
	}

	existFlag := false
	for _, v := range validCodes {
		if code == v {
			existFlag = true
			break
		}
	}

	if !existFlag {
		e := NewError(ET_ParamsError, ED_InvalidTotpKey).WithParam("code", code)
		mlog.Log(e.String())
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
