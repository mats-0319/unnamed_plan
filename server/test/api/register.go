package api

import (
	"fmt"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func Register() {
	testCase("duplicate register", registerCase_Duplicate)
	testCase("success", registerCase_Success)
}

func registerCase_Duplicate() string {
	res := httpInvoke(api.URI_Register, `{"user_name":"admin","password":""}`)
	if res.IsSuccess || res.Err != utils.ErrUserExist().Error() {
		return unknownError
	}

	return ""
}

func registerCase_Success() string {
	pwd := utils.CalcSHA256("123456")

	res := httpInvoke(api.URI_Register, fmt.Sprintf(`{"user_name":"new_user","password":"%s"}`, pwd))
	if !res.IsSuccess {
		return res.Err
	}

	count, data, err := dal.ListUser(api.Pagination{Size: 10, Num: 1})
	if count != 4 || len(data) != 4 || err != nil {
		return unknownError
	}

	return ""
}
