package handlers

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/cmd/model"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func Register(ctx *mhttp.Context) {
	req := &api.RegisterReq{}
	if !ctx.ParseParams(req) {
		return
	}

	pwd := utils.GeneratePwdHash(req.Password)

	user := &model.User{
		UserName: req.UserName,
		Nickname: req.UserName,
		Password: pwd,
	}
	err := dal.CreateUser(user)
	if err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.RegisterRes{}
}
