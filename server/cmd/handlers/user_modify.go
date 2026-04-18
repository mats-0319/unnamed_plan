package handlers

import (
	"encoding/base32"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/password"
)

func ModifyUser(ctx *mhttp.Context) {
	req := &api.ModifyUserReq{}
	if !ctx.ParseParams(req) {
		return
	}

	operator, e := dal.GetUser(ctx.UserName)
	if e != nil {
		ctx.ResData = e
		return
	}

	modifyNicknameFlag := len(req.Nickname) > 0 && req.Nickname != operator.Nickname
	modifyTOTPKeyFlag := len(req.TOTPKey) > 0 && req.TOTPKey != operator.TOTPKey
	if !modifyNicknameFlag && len(req.Password) < 1 && req.EnableMFA == operator.EnableMFA && !modifyTOTPKeyFlag {
		e := utils.ErrNoChanges().WithParam("operator", operator.UserName)
		ctx.ResData = e
		mlog.Error(e.String())
		return
	}

	if modifyNicknameFlag {
		operator.Nickname = req.Nickname
	}
	if len(req.Password) > 0 {
		if e := password.VerifyPassword(req.Password, operator.Password); e == nil { // in modify, same pwd is wrong
			e := utils.ErrSamePassword()
			ctx.ResData = e
			mlog.Error(e.String())
			return
		}

		operator.Password = password.GeneratePassword(req.Password)
	}
	if modifyTOTPKeyFlag {
		if e := isValidTOTPKey(req.TOTPKey); e != nil {
			ctx.ResData = e
			return
		}

		operator.TOTPKey = req.TOTPKey
	}
	if req.EnableMFA != operator.EnableMFA {
		operator.EnableMFA = req.EnableMFA // 先判断totp key，到这里就不用重复判断其有效性了
	}

	if e := dal.UpdateUser(operator); e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.ModifyUserRes{}
}

func isValidTOTPKey(key string) *utils.Error {
	bytes, err := base32.StdEncoding.DecodeString(key)
	if err != nil || !(0 < len(bytes) && len(bytes) <= 10) { // 空字符串是有效的base32字符串，但不应该是有效的key
		e := utils.ErrInvalidTOTPKey().WithParam("totp key", key).WithCause(err)
		mlog.Error(e.String())
		return e
	}

	return nil
}
