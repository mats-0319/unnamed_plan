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

	if len(ctx.UserName) < 1 {
		e := utils.ErrInvalidParams().WithParam("operator", ctx.UserName)
		mlog.Error(e.String())
		ctx.ResData = e
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
	if req.EnableMFA || modifyTOTPKeyFlag { // 准备启用MFA/更换新的totp key时，检查totp key是否有效
		if e := isValidTOTPKey(req.TOTPKey); e != nil {
			ctx.ResData = e
			return
		}
	}

	// modify
	if modifyNicknameFlag {
		operator.Nickname = req.Nickname
	}
	if len(req.Password) > 0 {
		if password.VerifyPassword(req.Password, operator.Password) == nil { // in modify, same pwd is wrong
			e := utils.ErrSamePassword()
			ctx.ResData = e
			mlog.Error(e.String())
			return
		}

		operator.Password = password.GeneratePassword(req.Password)
	}
	if req.EnableMFA != operator.EnableMFA {
		operator.EnableMFA = req.EnableMFA
	}
	if modifyTOTPKeyFlag {
		operator.TOTPKey = req.TOTPKey
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
