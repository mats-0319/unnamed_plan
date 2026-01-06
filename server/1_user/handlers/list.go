package handlers

import (
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ListUser(ctx *mhttp.Context) {
	req := &api.ListUserReq{}
	if !ctx.ParseParams(req) {
		return
	}

	operator, err := dal.GetUser(ctx.UserID)
	if err != nil {
		ctx.ResData = err
		return
	}

	if !operator.IsAdmin {
		e := NewError(ET_OperatorError, ED_NeedAdmin).WithParam("operator", operator.UserName)
		ctx.ResData = e
		mlog.Log(e.String())
		return
	}

	count, users, err := dal.ListUsers(req.Page)
	if err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.ListUserRes{
		ResBase: api.ResBaseSuccess,
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
			UserName:  v.UserName,
			Nickname:  v.Nickname,
			TotpKey:   v.TotpKey,
			IsAdmin:   v.IsAdmin,
			LastLogin: v.LastLogin,
		}
	}

	return res
}
