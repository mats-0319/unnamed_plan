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
		AccessToken: r.Header.Get(HTTPHeader_AccessToken),
	}
}

// ParseParams 这个函数读取了http req.body，这一结构被限制只能读取一次，
// 所以包括请求处理函数和中间件在内，如果有多个函数均调用该函数，会出现错误。
func (ctx *Context) ParseParams(obj any) bool {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		e := ErrServerInternalError().WithCause(err)
		ctx.ResData = e
		mlog.Error(e.String())
		return false
	}

	if err := json.Unmarshal(bodyBytes, obj); err != nil {
		e := ErrDeserializeHTTPReqParam().WithCause(err)
		ctx.ResData = e
		mlog.Error(e.String())
		return false
	}

	return true
}
