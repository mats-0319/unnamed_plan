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
	modifyTotpKeyFlag := len(req.TotpKey) > 0 && req.TotpKey != operator.TotpKey
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
		if e := password.VerifyPassword(req.Password, operator.Password); e == nil { // in modify, same pwd is wrong
			e := utils.ErrSamePassword()
			ctx.ResData = e
			mlog.Error(e.String())
			return
		}

		operator.Password = password.GeneratePassword(req.Password)
	}
	if req.Enable2FA != operator.Enable2FA {
		if req.Enable2FA {
			if e := isValidTotpKey(req.TotpKey); e != nil { // 想要启用2FA，但是totp key无效
				ctx.ResData = e
				return
			}
		}

		operator.Enable2FA = req.Enable2FA
	}
	if modifyTotpKeyFlag {
		if e := isValidTotpKey(req.TotpKey); e != nil {
			ctx.ResData = e
			return
		}

		operator.TotpKey = req.TotpKey
	}

	if e := dal.UpdateUser(operator); e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.ModifyUserRes{}
}

func isValidTotpKey(key string) *utils.Error {
	bytes, err := base32.StdEncoding.DecodeString(key)
	if err != nil || !(0 < len(bytes) && len(bytes) <= 10) { // 空字符串是有效的base32字符串，但不应该是有效的key
		e := utils.ErrInvalidTotpKey().WithParam("totp key", key).WithCause(err)
		mlog.Error(e.String())
		return e
	}

	return nil
}
