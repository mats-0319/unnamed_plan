package mhttp

import (
	"encoding/json"
	"io"
	"net/http"

	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	ResData any // allow: errStr/error/struct

	noResponse bool // forward req will use 'io.copy' instead of 'ctx.response'
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
	}
}

func (ctx *Context) Bind(obj any) error {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		mlog.Log("read req body failed", mlog.Field("error", err))
		return err
	}

	err = json.Unmarshal(bodyBytes, obj)
	if err != nil {
		mlog.Log("parse req body failed", mlog.Field("error", err))
		return err
	}

	return nil
}

func (ctx *Context) response() {
	if ctx.noResponse {
		return
	}

	var obj any
	switch v := ctx.ResData.(type) {
	case string:
		if len(v) > 0 {
			obj = &api.ResBase{Err: v}
		}
	case error:
		obj = &api.ResBase{Err: v.Error()}
	default:
		obj = v
	}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		mlog.Log("serialize handlers res to json failed", mlog.Field("error", err))
	}

	_, err = ctx.Writer.Write(jsonBytes)
	if err != nil {
		mlog.Log("handlers res failed", mlog.Field("error", err))
		return
	}
}

func (ctx *Context) Forward(newHost string) error {
	url := "http://" + newHost + ctx.Request.RequestURI
	req, err := http.NewRequest("POST", url, ctx.Request.Body)
	if err != nil {
		mlog.Log("new req failed", mlog.Field("error", err))
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		mlog.Log("handlers invoke failed", mlog.Field("error", err))
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	_, err = io.Copy(ctx.Writer, res.Body)
	if err != nil {
		mlog.Log("return res to web failed", mlog.Field("error", err))
		return err
	}

	ctx.noResponse = true

	return nil
}
