package api

import (
	"fmt"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func Login() {
	testCase("user not exist", loginCase_UserNotExist)
	testCase("wrong pwd", loginCase_WrongPwd)
	testCase("wrong totp code", loginCase_WrongTotpCode)
	testCase("success", loginCase_Success)
}

func loginCase_UserNotExist() string {
	res := httpInvoke(api.URI_Login, `{"user_name":"not exist","password":"","totp_code":""}`, "")
	if res.IsSuccess || res.Err != utils.ErrUserNotFound().Error() {
		return unknownError
	}

	return ""
}

func loginCase_WrongPwd() string {
	res := httpInvoke(api.URI_Login, `{"user_name":"admin","password":"wrong pwd","totp_code":""}`, "")
	if res.IsSuccess || res.Err != utils.ErrWrongPassword().Error() {
		return unknownError
	}

	return ""
}

func loginCase_WrongTotpCode() string {
	res := httpInvoke(api.URI_Login, fmt.Sprintf(`{"user_name":"user_with_totp","password":"%s","totp_code":"000000"}`, pwd), "")
	if res.IsSuccess || res.Err != utils.ErrWrongTotpCode().Error() {
		return unknownError
	}

	return ""
}

func loginCase_Success() string {
	res := httpInvoke(api.URI_Login, fmt.Sprintf(`{"user_name":"user","password":"%s","totp_code":""}`, pwd), "")
	if !res.IsSuccess {
		return res.Err
	}

	return ""
}
