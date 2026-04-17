package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	"github.com/mats0319/unnamed_plan/server/internal/http/middleware"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"github.com/mats0319/unnamed_plan/server/internal/utils/password"
)

func Login(ctx *mhttp.Context) {
	req := &api.LoginReq{}
	if !ctx.ParseParams(req) {
		return
	}

	if len(req.UserName) < 1 || len(req.Password) < 1 {
		e := utils.ErrInvalidParams().WithParam("user name", req.UserName).WithParam("password", req.Password)
		mlog.Error(e.String())
		ctx.ResData = e
		return
	}

	user, e := dal.GetUser(req.UserName)
	if e != nil {
		ctx.ResData = e
		return
	}

	if e := password.VerifyPassword(req.Password, user.Password); e != nil {
		ctx.ResData = e
		return
	}

	if user.Enable2FA {
		ctx.ResData = &api.LoginRes{Enable2FA: true, MfaToken: middleware.GenerateMfaToken(user.UserName)}

		return
	}

	_ = dal.UpdateUser(user) // modify user.UpdatedAt

	ctx.Writer.Header().Set(utils.HttpHeader_AccessToken, middleware.GenerateApiAccessToken(user.UserName))

	ctx.ResData = &api.LoginRes{
		UserName: user.UserName,
		Nickname: user.Nickname,
		IsAdmin:  user.IsAdmin,
	}
}
