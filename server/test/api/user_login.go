package api

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"math"
	"time"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func Login() {
	testCase("user not exist", loginCase_UserNotExist)
	testCase("wrong pwd", loginCase_WrongPwd)
	testCase("success - disable 2fa", loginCase_SuccessDisable2fa)
	testCase("success - enable 2fa", loginCase_SuccessEnable2fa)
}

func loginCase_UserNotExist() string {
	res := httpInvoke(api.URI_Login, `{"user_name":"not exist","password":"123"}`, "")
	if res.IsSuccess || res.Err != utils.ErrUserNotFound().Error() {
		return unknownError
	}

	return ""
}

func loginCase_WrongPwd() string {
	res := httpInvoke(api.URI_Login, `{"user_name":"admin","password":"wrong pwd"}`, "")
	if res.IsSuccess || res.Err != utils.ErrWrongPassword().Error() {
		return unknownError
	}

	return ""
}

func loginCase_SuccessDisable2fa() string {
	res := httpInvoke(api.URI_Login, fmt.Sprintf(`{"user_name":"user","password":"%s"}`, pwdSHA256), "")
	if !res.IsSuccess {
		return res.Err
	}

	return ""
}

func loginCase_SuccessEnable2fa() string {
	res := httpInvoke(api.URI_Login, fmt.Sprintf(`{"user_name":"user_with_totp","password":"%s"}`, pwdSHA256), "")
	if !res.IsSuccess {
		return res.Err
	}

	mfaToken := res.Data.MfaToken
	if len(mfaToken) < 1 {
		return unknownError
	}

	totpCode := calcTotpCode([]byte("mario"), iTob(time.Now().Unix()/30))
	res = httpInvoke(api.URI_LoginTotp, fmt.Sprintf(`{"mfa_token":"%s","totp_code":"%s"}`, mfaToken, totpCode), "")
	if !res.IsSuccess {
		return res.Err
	}

	return ""
}

// copy from handler
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
	byteSlice := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		byteSlice[i] = byte(integer & 0xff)
		integer = integer >> 8
	}
	return byteSlice
}
