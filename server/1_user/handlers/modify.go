package handlers

import (
	"encoding/base32"

	. "github.com/mats0319/unnamed_plan/server/internal/const"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ModifyUser(ctx *mhttp.Context) {
	req := &api.ModifyUserReq{}
	if !ctx.ParseParams(req) {
		return
	}

	operator, err := dal.GetUser(ctx.UserID)
	if err != nil {
		ctx.ResData = err
		return
	}

	if len(req.Nickname) > 0 {
		operator.Nickname = req.Nickname
	}
	if len(req.Password) > 0 {
		newPassword := utils.CalcSHA256(req.Password, operator.Salt)
		if newPassword == operator.Password {
			e := NewError(ET_ParamsError, ED_SamePwd)
			mlog.Log(e.String())
			ctx.ResData = e
			return
		}

		operator.Password = newPassword
	}
	if req.ModifyTkFlag {
		bytes, err := base32.StdEncoding.DecodeString(req.TotpKey)
		if err != nil || len(bytes) > 10 {
			e := NewError(ET_ParamsError, ED_InvalidTotpKey).WithParam("totp key", req.TotpKey).WithCause(err)
			mlog.Log(e.String())
			ctx.ResData = e
			return
		}

		operator.TotpKey = req.TotpKey
	}

	err = dal.UpdateUser(operator)
	if err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.ModifyUserRes{
		ResBase: api.ResBase{IsSuccess: true},
	}
}
