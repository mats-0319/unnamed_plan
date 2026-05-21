package handlers

import (
	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/cmd/handlers/mfa"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func SetMFAStatus(ctx *mhttp.Context) {
	req := &api.SetMFAStatusReq{}
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

	if !req.EnableMFA { // 禁用MFA
		if operator.EnableMFA { // 仅在启用状态下更新数据库，如已是禁用状态，直接返回
			operator.EnableMFA = false
			if e := dal.UpdateUser(operator); e != nil {
				ctx.ResData = e
				return
			}
		}

		ctx.ResData = &api.SetMFAStatusRes{}

		return
	}

	/* 启用MFA */

	if req.ApplyNewKeyFlag { // 申请了新的totp key
		totpKey, e := mfa.VerifyTOTPKey(ctx.UserName, req.TOTPCode)
		if e != nil {
			ctx.ResData = e
			return
		}
		operator.TOTPKey = totpKey
	} else { // 使用历史totp key
		if e := mfa.VerifyTOTPCode(req.TOTPCode, operator.TOTPKey); e != nil {
			ctx.ResData = e
			return
		}
	}

	operator.EnableMFA = true
	if e := dal.UpdateUser(operator); e != nil {
		ctx.ResData = e
		return
	}

	ctx.ResData = &api.SetMFAStatusRes{}
}
