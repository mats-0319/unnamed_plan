package mhttp

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/mats0319/unnamed_plan/server/internal/const"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Origin  string // when forward, set http req header 'Origin'

	UserID      uint
	AccessToken string // 登录成功获得，后续请求均需要在请求头带上该参数

	ResHeaders map[string]string // header - value
	ResData    any               // allow: errStr/error/struct/[]byte
}

func NewContext(w http.ResponseWriter, r *http.Request, origin string) *Context {
	userIDStr := r.Header.Get(mconst.HttpHeader_UserID)
	userID, _ := strconv.Atoi(userIDStr)
	token := r.Header.Get(mconst.HttpHeader_AccessToken)

	return &Context{
		Writer:      w,
		Request:     r,
		Origin:      origin,
		UserID:      uint(userID),
		AccessToken: token,
		ResHeaders:  make(map[string]string),
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
	for header, value := range ctx.ResHeaders {
		ctx.Writer.Header().Set(header, value)
	}

	code, resBytes := serializeRes(ctx.ResData)

	if code != http.StatusOK {
		ctx.Writer.WriteHeader(code)
	}

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
		obj = &api.ResBase{Err: v.Error()}

		switch v.Typ {
		case ET_ServerInternalError:
			code = http.StatusInternalServerError
		case ET_UnauthorizedError:
			code = http.StatusUnauthorized
		}
	case []byte: // forward res, no marshal
		return code, v
	default: // struct
	}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		// 因为这里已经给resBytes定型了，所以返回错误也没啥能做的，就不返回了
		mlog.Log("serialize handlers res to json failed", mlog.Field("error", err))
		return code, nil
	}

	return code, jsonBytes
}
