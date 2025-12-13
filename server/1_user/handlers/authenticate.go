package handlers

import (
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
)

func Authenticate(ctx *mhttp.Context) {
	// 在中间件验证user id和access token，能到这说明验证通过了，直接返回成功就行
	ctx.ResData = &api.AuthenticateRes{
		ResBase: api.ResBase{IsSuccess: true},
	}
}
