package handlers

import (
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
)

func Authenticate(ctx *mhttp.Context) {
	user, err := dal.GetUser(ctx.UserID)
	if err != nil {
		ctx.ResData = err
		return
	}

	ctx.ResData = &api.AuthenticateRes{
		ResBase:  api.ResBaseSuccess,
		UserID:   user.ID,
		UserName: user.UserName,
		Nickname: user.Nickname,
		IsAdmin:  user.IsAdmin,
	}
}
