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

	AccessToken string // 登录成功获得，后续请求均需要在请求头带上该参数
	UserName    string // user name

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

	err = json.Unmarshal(bodyBytes, obj)
	if err != nil {
		e := ErrJsonUnmarshal().WithCause(err)
		mlog.Error(e.String())
		ctx.ResData = e
		return false
	}

	return true
}
