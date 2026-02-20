package handlers

import (
	"encoding/base32"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ModifyUser(ctx *mhttp.Context) {
	req := &api.ModifyUserReq{}
	if !ctx.ParseParams(req) {
		return
	}

	operator, err := dal.GetUser(ctx.UserName)
	if err != nil {
		ctx.ResData = err
		return
	}

	modifyNicknameFlag := len(req.Nickname) > 0 && req.Nickname != operator.Nickname
	modifyTotpKeyFlag := len(req.TotpKey) > 0 && req.TotpKey != operator.TotpKey && len(req.TotpKey) > 0
	if !modifyNicknameFlag && len(req.Password) < 1 && req.Enable2FA == operator.Enable2FA && !modifyTotpKeyFlag {
		e := utils.ErrNoChanges().WithParam("operator", operator.UserName)
		ctx.ResData = e
		mlog.Error(e.String())
		return
	}

	if modifyNicknameFlag {
		operator.Nickname = req.Nickname
	}
	if len(req.Password) > 0 {
		err = utils.VerifyPassword(req.Password, operator.Password) // in modify, same pwd is invalid
		if err == nil {
			e := utils.ErrSamePwd()
			ctx.ResData = e
			mlog.Error(e.String())
			return
		}

		operator.Password = utils.GeneratePwdHash(req.Password)
	}
	err = isValidTotpKey(req.TotpKey)
	if req.Enable2FA != operator.Enable2FA {
		if req.Enable2FA && err != nil { // 想要启用2FA，但是totp key无效
			ctx.ResData = err
			return
		}

		operator.Enable2FA = req.Enable2FA
	}
	if modifyTotpKeyFlag {
		if err != nil {
			ctx.ResData = err
			return
		}

		operator.TotpKey = req.TotpKey
	}

	err = dal.UpdateUser(operator)
	if err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.ModifyUserRes{}
}

func isValidTotpKey(key string) *utils.Error {
	bytes, err := base32.StdEncoding.DecodeString(key)
	if len(key) < 1 || err != nil || len(bytes) > 10 { // 空字符串是有效的base32字符串，但是不应该是有效的key
		e := utils.ErrInvalidTotpKey().WithParam("totp key", key).WithCause(err)
		mlog.Error(e.String())
		return e
	}

	return nil
}
