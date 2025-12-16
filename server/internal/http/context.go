package mhttp

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	. "github.com/mats0319/unnamed_plan/server/internal/const"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	UserID      uint
	AccessToken string // 登录成功获得，后续请求均需要在请求头带上该参数

	ResData    any               // allow: errStr/error/struct/[]byte
	ResHeaders map[string]string // header - value
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	userIDStr := r.Header.Get(HttpHeader_UserID)
	userID, _ := strconv.Atoi(userIDStr)
	token := r.Header.Get(HttpHeader_AccessToken)

	return &Context{
		Writer:      w,
		Request:     r,
		UserID:      uint(userID),
		AccessToken: token,
		ResHeaders:  make(map[string]string),
	}
}

func (ctx *Context) response() {
	for header, value := range ctx.ResHeaders {
		ctx.Writer.Header().Set(header, value)
	}

	resBytes, err := serializeRes(ctx.ResData)
	if err != nil {
		mlog.Log("serialize res failed", mlog.Field("error", err))
		return
	}

	_, err = ctx.Writer.Write(resBytes)
	if err != nil {
		mlog.Log("response failed", mlog.Field("error", err))
		return
	}
}

func (ctx *Context) ParseParams(obj any) bool {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		e := NewError(ET_ServerInternalError, ED_IORead).WithCause(err)
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

func (ctx *Context) Forward(url string, r io.Reader) {
	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		e := NewError(ET_ServerInternalError, ED_InvalidHttpRequest).WithCause(err)
		mlog.Log(e.String())
		ctx.ResData = e
		return
	}

	for _, header := range HttpHeaderList { // forward our own http header
		value := ctx.Request.Header.Get(header)
		req.Header.Add(header, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		e := NewError(ET_ServerInternalError, ED_HttpInvoke).WithCause(err)
		mlog.Log(e.String())
		return
	}
	defer func() {
		_ = res.Body.Close()
	}()

	headers := make(map[string]string)
	for header, value := range res.Header {
		for _, v := range HttpHeaderList { // 记录我们自定义的'header'
			if header == v {
				headers[header] = value[0]
			}
		}
	}

	// 这里我们希望统一使用ctx.response设置响应头和返回值，所以不使用io.copy直接复制res.body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		e := NewError(ET_ServerInternalError, ED_IORead).WithCause(err)
		mlog.Log(e.String())
		return
	}

	ctx.ResHeaders = headers
	ctx.ResData = bodyBytes

	return
}

func serializeRes(obj any) ([]byte, error) {
	switch v := obj.(type) {
	case string: // err str
		if len(v) > 0 {
			obj = &api.ResBase{Err: v}
		}
	case error:
		obj = &api.ResBase{Err: v.Error()}
	case []byte: // forward res, no marshal
		return v, nil
	default: // struct
	}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		mlog.Log("serialize handlers res to json failed", mlog.Field("error", err))
		return nil, err
	}

	return jsonBytes, nil
}
