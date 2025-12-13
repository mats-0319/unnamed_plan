package handlers

import (
	"github.com/mats0319/unnamed_plan/server/1_user/db"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func CreateUser(ctx *mhttp.Context) {
	req := &api.CreateUserReq{}
	if err := ctx.ParseParams(req); err != nil {
		mlog.Log("parse params failed", mlog.Field("error", err))
		ctx.ResData = err
		return
	}

	salt := string(utils.GenerateRandomBytes(10))
	pwd := utils.CalcSHA256(req.Password + salt)

	user := &model.User{
		Name:     req.UserName,
		Nickname: req.UserName,
		Password: pwd,
		Salt:     salt,
	}

	err := db.CreateUser(user)
	if err != nil {
		// todo: distinguish db error and 'user name' exist error
		mlog.Log("create user failed", mlog.Field("error", err))
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.CreateUserRes{
		ResBase: api.ResBase{IsSuccess: true},
	}
}
