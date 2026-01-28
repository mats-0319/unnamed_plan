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
	UserID      uint

	ResData any // expect: *utils.Error / *api.resStruct
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	token := r.Header.Get(HttpHeader_AccessToken)

	return &Context{
		Writer:      w,
		Request:     r,
		AccessToken: token,
	}
}

func (ctx *Context) ParseParams(obj any, r ...io.Reader) bool {
	var reader io.Reader
	if len(r) > 0 {
		reader = r[0]
	} else {
		reader = ctx.Request.Body
	}

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		e := NewError(ET_ServerInternalError).WithCause(err)
		mlog.Log(e.String())
		ctx.ResData = e
		return false
	}

	err = json.Unmarshal(bodyBytes, obj)
	if err != nil {
		e := NewError(ET_ParamsError, ED_JsonUnmarshal).WithCause(err)
		mlog.Log(e.String())
		ctx.ResData = e
		return false
	}

	return true
}
