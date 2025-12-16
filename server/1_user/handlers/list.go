package handlers

import (
	"github.com/mats0319/unnamed_plan/server/1_user/db"
	. "github.com/mats0319/unnamed_plan/server/internal/const"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

func ListUser(ctx *mhttp.Context) {
	req := &api.ListUserReq{}
	if !ctx.ParseParams(req) {
		return
	}

	operator, err := db.GetUser(ctx.UserID)
	if err != nil {
		ctx.ResData = err
		return
	}

	if !operator.IsAdmin {
		e := NewError(ET_ParamsError, ED_NeedAdmin).WithParam("operator", operator.UserName)
		mlog.Log(e.String())
		ctx.ResData = e
		return
	}

	count, users, err := db.ListUsers(req.Page)
	if err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.ListUserRes{
		ResBase: api.ResBase{IsSuccess: true},
		Amount:  count,
		Page:    req.Page,
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
			Name:      v.UserName,
			Nickname:  v.Nickname,
			TotpKey:   v.TotpKey,
			LastLogin: v.LastLogin,
		}
	}

	return res
}
