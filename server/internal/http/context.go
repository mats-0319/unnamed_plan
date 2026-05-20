package mhttp

import (
	"encoding/json"
	"io"
	"net/http"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

type Context struct {
	writer  http.ResponseWriter
	request *http.Request

	AccessToken string // 登录成功后签发，需要身份认证的接口应包含此项
	UserName    string // parse from 'AccessToken'

	ResData any // expect: *utils.Error / *api.resStruct
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		writer:      w,
		request:     r,
		AccessToken: r.Header.Get(utils.HTTPHeader_AccessToken),
	}
}

// ParseParams 这个函数读取了http req.body，这一结构被限制只能读取一次，
// 所以包括请求处理函数和中间件在内，如果有多个函数均调用该函数，会出现错误。
func (ctx *Context) ParseParams(obj any) bool {
	bodyBytes, err := io.ReadAll(ctx.request.Body)
	if err != nil {
		e := utils.ErrServerInternalError().WithCause(err)
		mlog.Error(e.String())
		ctx.ResData = e
		return false
	}

	if err := json.Unmarshal(bodyBytes, obj); err != nil {
		e := utils.ErrDeserializeReqParam().WithCause(err)
		mlog.Error(e.String())
		ctx.ResData = e
		return false
	}

	return true
}

func (ctx *Context) SetHeader(key string, value string) {
	ctx.writer.Header().Set(key, value)
}
