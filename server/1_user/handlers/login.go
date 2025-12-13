package handlers

import (
	"fmt"
	"time"

	"github.com/mats0319/unnamed_plan/server/1_user/db"
	"github.com/mats0319/unnamed_plan/server/1_user/middleware"
	mconst "github.com/mats0319/unnamed_plan/server/internal/const"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func Login(ctx *mhttp.Context) {
	req := &api.LoginReq{}
	if err := ctx.ParseParams(req); err != nil {
		mlog.Log("parse params failed", mlog.Field("error", err))
		ctx.ResData = err
		return
	}

	user, err := db.GetUser(req.UserName)
	if err != nil {
		mlog.Log("get user failed", mlog.Field("error", err))
		ctx.ResData = err
		return
	}

	if utils.CalcSHA256(req.Password, user.Salt) != user.Password {
		str := fmt.Sprintf("invalid username or password: %s, %s", req.UserName, req.Password)
		mlog.Log("login failed", str)
		ctx.ResData = str
		return
	}

	if len(user.TotpKey) > 0 {
		// check totp code
	}

	// todo: token use handlers head

	user.LastLogin = time.Now().UnixMilli()
	if err = db.UpdateUser(user); err != nil {
		mlog.Log("update user failed", mlog.Field("error", err))
		ctx.ResData = err
		return
	}

	token := string(utils.GenerateRandomBytes(10))
	middleware.SetToken(user.ID, token)

	ctx.ResHeaders[mconst.HttpHeader_AccessToken] = token

	ctx.ResData = &api.LoginRes{
		ResBase:  api.ResBase{IsSuccess: true},
		UserID:   user.ID,
		UserName: user.Name,
		Nickname: user.Nickname,
		IsAdmin:  user.IsAdmin,
	}
}
