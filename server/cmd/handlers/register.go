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

	salt := utils.GenerateRandomBytes[string](10)
	pwd := utils.HmacSHA256(req.Password, salt)

	user := &model.User{
		UserName: req.UserName,
		Nickname: req.UserName,
		Password: pwd,
		Salt:     salt,
	}
	err := dal.CreateUser(user)
	if err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.RegisterRes{}
}
