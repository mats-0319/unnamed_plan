package handlers

import (
	"github.com/mats0319/unnamed_plan/server/1_user/db"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func CreateUser(ctx *mhttp.Context) {
	req := &api.CreateUserReq{}
	if !ctx.ParseParams(req) {
		return
	}

	salt := string(utils.GenerateRandomBytes(10))
	pwd := utils.CalcSHA256(req.Password + salt)

	user := &model.User{
		UserName: req.UserName,
		Nickname: req.UserName,
		Password: pwd,
		Salt:     salt,
	}

	err := db.CreateUser(user)
	if err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.CreateUserRes{
		ResBase: api.ResBase{IsSuccess: true},
	}
}
