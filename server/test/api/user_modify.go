package api

import (
	"fmt"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ModifyUser() {
	testCase("no changes", modifyUserCase_NoChanges)
	testCase("same pwd", modifyUserCase_SamePwd)
	testCase("invalid totp key", modifyUserCase_InvalidTotpKey)
	testCase("success", modifyUserCase_Success)
}

func modifyUserCase_NoChanges() string {
	res := httpInvoke(api.URI_ModifyUser, `{"nickname":"","password":"","enable_2fa":false,"totp_key":""}`, accessToken_User)
	if res.IsSuccess || res.Err != utils.ErrNoChanges().Error() {
		return unknownError
	}

	return ""
}

func modifyUserCase_SamePwd() string {
	res := httpInvoke(api.URI_ModifyUser, fmt.Sprintf(`{"nickname":"","password":"%s","enable_2fa":false,"totp_key":""}`, pwd), accessToken_User)
	if res.IsSuccess || res.Err != utils.ErrSamePassword().Error() {
		return unknownError
	}

	return ""
}

func modifyUserCase_InvalidTotpKey() string {
	res := httpInvoke(api.URI_ModifyUser, `{"nickname":"","password":"","enable_2fa":true,"totp_key":"123"}`, accessToken_User)
	if res.IsSuccess || res.Err != utils.ErrInvalidTotpKey().Error() {
		return unknownError
	}

	return ""
}

func modifyUserCase_Success() string {
	res := httpInvoke(api.URI_ModifyUser, `{"nickname":"123","password":"","enable_2fa":false,"totp_key":""}`, accessToken_User)
	if !res.IsSuccess {
		return res.Err
	}

	user, err := dal.GetUser("user")
	if (user != nil && user.Nickname != "123") || err != nil {
		return unknownError
	}

	return ""
}
