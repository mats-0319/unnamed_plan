package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func Login(ctx *mhttp.Context) {
	req := &api.LoginReq{}

	if err := ctx.Bind(req); err != nil {
		mlog.Log("bind params failed", mlog.Field("error", err))
		ctx.ResData = err
		return
	}

	user, err := getUserByName(req.UserName)
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

	ctx.ResData = &api.LoginRes{
		ResBase:    api.ResBase{IsSuccess: true},
		Nickname:   user.Nickname,
		Permission: user.Permission,
	}
}

func getUserByName(userName string) (*model.User, error) {
	qu := dal.Q.User
	res, err := qu.WithContext(context.TODO()).Where(qu.Name.Eq(userName)).Find()
	if err != nil {
		mlog.Log("query user failed", mlog.Field("error", err))
		return nil, err
	}

	if len(res) < 1 {
		str := "no record for user name: " + userName
		mlog.Log("query user failed", str)
		return nil, errors.New(str)
	}

	return res[0], nil
}
