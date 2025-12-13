package handlers

import (
	"github.com/mats0319/unnamed_plan/server/1_user/db"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ModifyUser(ctx *mhttp.Context) {
	req := &api.ModifyUserReq{}
	if err := ctx.ParseParams(req); err != nil {
		mlog.Log("parse params failed", mlog.Field("error", err))
		ctx.ResData = err
		return
	}

	operator, err := db.GetUser(ctx.UserID)
	if err != nil {
		mlog.Log("get user failed", mlog.Field("error", err))
		ctx.ResData = err
		return
	}

	if len(req.Nickname) > 0 {
		operator.Nickname = req.Nickname
	}
	if len(req.Password) > 0 {
		if utils.CalcSHA256(req.Password+operator.Salt) == operator.Password {
			mlog.Log("new password can't be same with the old")
			ctx.ResData = "try to reset same password"
			return
		}

		operator.Password = req.Password
	}
	if req.ModifyTkFlag {
		operator.TotpKey = req.TotpKey
	}

	err = db.UpdateUser(operator)
	if err != nil {
		mlog.Log("update user failed", mlog.Field("error", err))
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.ModifyUserRes{
		ResBase: api.ResBase{IsSuccess: true},
	}
}
