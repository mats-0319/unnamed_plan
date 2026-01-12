package handlers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"time"

	mconst "github.com/mats0319/unnamed_plan/server/internal/const"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/http/middleware"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func Login(ctx *mhttp.Context) {
	req := &api.LoginReq{}
	if !ctx.ParseParams(req) {
		return
	}

	user, err := dal.GetUser(req.UserName)
	if err != nil {
		ctx.ResData = err
		return
	}

	if utils.CalcSHA256(req.Password, []byte(user.Salt)...) != user.Password {
		e := utils.NewError(utils.ET_ParamsError, utils.ED_InvalidPwd).WithParam("user name", req.UserName)
		ctx.ResData = e
		mlog.Log(e.String())
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

	err = dal.UpdateUser(user)
	if err != nil {
		ctx.ResData = err
		return
	}

	// token
	ctx.Writer.Header().Set(mconst.HttpHeader_AccessToken, middleware.GenToken(user.ID))

	ctx.ResData = &api.LoginRes{
		UserID:   user.ID,
		UserName: user.UserName,
		Nickname: user.Nickname,
		IsAdmin:  user.IsAdmin,
	}
}

// verifyTotpCode totpKey is base32 encoded
func verifyTotpCode(code string, totpKey string) *utils.Error {
	if len(code) != 6 {
		e := utils.NewError(utils.ET_ParamsError, utils.ED_InvalidTotpCode).WithParam("code", code)
		mlog.Log(e.String())
		return e
	}

	key := make([]byte, 10)
	n, err := base32.StdEncoding.Decode(key, []byte(totpKey))
	if err != nil {
		e := utils.NewError(utils.ET_ServerInternalError).WithCause(err)
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
		e := utils.NewError(utils.ET_ParamsError, utils.ED_InvalidTotpCode).WithParam("code", code)
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
