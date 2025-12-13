package handlers

import (
	"github.com/mats0319/unnamed_plan/server/1_user/db"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

func ListUser(ctx *mhttp.Context) {
	req := &api.ListUserReq{}
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

	if !operator.IsAdmin {
		str := mlog.Field("invalid operator", operator.ID)
		mlog.Log("check failed", str)
		ctx.ResData = str
		return
	}

	count, users, err := db.ListUsers(req.Page)
	if err != nil {
		mlog.Log("list users failed", mlog.Field("error", err))
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.ListUserRes{
		ResBase: api.ResBase{IsSuccess: true},
		Amount:  count,
		Users:   usersFromDBToHttp(users),
	}
}

func usersFromDBToHttp(users []*model.User) []*api.User {
	res := make([]*api.User, len(users))
	for i, v := range users {
		res[i] = &api.User{
			ID:        v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Name:      v.Name,
			Nickname:  v.Nickname,
			TotpKey:   v.TotpKey,
			LastLogin: v.LastLogin,
		}
	}

	return res
}
