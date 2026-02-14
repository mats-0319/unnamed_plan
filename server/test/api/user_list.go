package api

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ListUser() {
	testCase("operator not admin", listUserCase_NotAdmin)
	testCase("success", listUserCase_Success)
}

func listUserCase_NotAdmin() string {
	loginCase_Success(false)() // set token

	res := httpInvoke(api.URI_ListUser, `{"page":{"size":10,"num":1}}`)
	if res.IsSuccess || res.Err != utils.ErrNeedAdmin().Error() {
		return unknownError
	}

	return ""
}

func listUserCase_Success() string {
	loginCase_Success(true)()

	res := httpInvoke(api.URI_ListUser, `{"page":{"size":10,"num":1}}`)
	if !res.IsSuccess {
		return res.Err
	}

	count, _, err := dal.ListUser(api.Pagination{Size: 10, Num: 1})
	if count != 4 || err != nil {
		return unknownError
	}

	return ""
}
