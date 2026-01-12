package mhttp

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/mats0319/unnamed_plan/server/internal/const"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
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
	token := r.Header.Get(mconst.HttpHeader_AccessToken)

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

// response 该函数不应该中途返回，一定要执行到write
func (ctx *Context) response() {
	code, resBytes := serializeRes(ctx.ResData)

	ctx.Writer.WriteHeader(code)

	_, err := ctx.Writer.Write(resBytes)
	if err != nil {
		mlog.Log("response failed", mlog.Field("error", err))
		return
	}
}

func serializeRes(obj any) (int, []byte) {
	code := http.StatusOK

	switch v := obj.(type) {
	case *Error:
		obj = &api.Response{Err: v.Error()}

		code = getHttpCode(v)
	default: // *api.resStruct(s)
		obj = &api.Response{
			IsSuccess: true,
			Data:      v,
		}
	}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		// 因为这里已经给resBytes定型了，返回错误也没啥能做的，就不返回了
		mlog.Log("serialize handlers res to json failed", mlog.Field("error", err))
		return code, nil
	}

	return code, jsonBytes
}

func getHttpCode(err *Error) int {
	code := http.StatusOK
	switch err.Typ {
	case ET_ServerInternalError:
		code = http.StatusInternalServerError
	case ET_UnauthorizedError:
		code = http.StatusUnauthorized
	}

	return code
}
