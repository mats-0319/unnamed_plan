package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	"github.com/mats0319/unnamed_plan/server/internal/utils/password"
)

func Register(ctx *mhttp.Context) {
	req := &api.RegisterReq{}
	if !ctx.ParseParams(req) {
		return
	}

	user := &model.User{
		UserName: req.UserName,
		Nickname: req.UserName,
		Password: password.GeneratePwdHash(req.Password),
	}

	if err := dal.CreateUser(user); err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.RegisterRes{}
}
