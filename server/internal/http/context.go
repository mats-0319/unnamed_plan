package mhttp

import (
	"encoding/json"
	"io"
	"net/http"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	AccessToken string // 登录成功后签发，需要身份认证的接口应包含此项
	UserName    string // parse from 'AccessToken'

	ResData any // expect: *utils.Error / *api.resStruct
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:      w,
		Request:     r,
		AccessToken: r.Header.Get(HttpHeader_AccessToken),
	}
}

func (ctx *Context) ParseParams(obj any) bool {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		e := ErrServerInternalError().WithCause(err)
		mlog.Error(e.String())
		ctx.ResData = e
		return false
	}

	if err := json.Unmarshal(bodyBytes, obj); err != nil {
		e := ErrJsonUnmarshal().WithCause(err)
		mlog.Error(e.String())
		ctx.ResData = e
		return false
	}

	return true
}
