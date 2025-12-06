package mhttp

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter

	ResData []byte

	IsSuccess   bool
	InvokeID    string
	InvokeChain []*InvokeItem
}

type InvokeItem struct {
	Timestamp int64
	URL       string
	Res       *api.ResBase
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Request:  r,
		Writer:   w,
		ResData:  []byte{},
		InvokeID: getInvokeID(r),
		InvokeChain: []*InvokeItem{{
			Timestamp: time.Now().UnixMilli(),
			URL:       r.URL.String(),
			Res:       nil,
		}},
	}
}

func (c *Context) Update(ctx *Context) {
	c.ResData = ctx.ResData
}

type ResBaseWrapper struct { // 包装一层，为了和其他api的res结构体在层级结构上一致
	Res *api.ResBase `json:"res"`
}

// NewResData 接受一个errStr字符串或者一个结构体参数，前者表示服务器内部错误，后者表示业务逻辑执行错误
func NewResData[T string | any](value T) []byte {
	var obj any

	switch v := any(value).(type) {
	case string:
		obj = &ResBaseWrapper{&api.ResBase{
			IsSuccess: false,
			Err:       v,
		}}
	default:
		obj = v
	}

	jsonBytes, _ := json.Marshal(obj)

	return jsonBytes
}

func response(writer http.ResponseWriter, data any) {
	h := writer.Header()
	if val := h["Content-Type"]; len(val) < 1 {
		h["Content-Type"] = []string{"application/json; charset=utf-8"}
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		mlog.Log("serialize http res to json failed", mlog.Field("error", err))
		return
	}

	_, err = writer.Write(jsonBytes)
	if err != nil {
		mlog.Log("http res failed", mlog.Field("error", err))
		return
	}
}

func getInvokeID(r *http.Request) string {
	formID := r.FormValue("invoke_id")
	if len(formID) > 0 {
		return formID
	}

	return uuid.NewString()
}
