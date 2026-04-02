package mhttp

import (
	"encoding/json"
	"net/http"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

type Response struct {
	IsSuccess bool   `json:"is_success"`
	Err       string `json:"err"`
	Data      any    `json:"data"`
}

// response 该函数不应该中途返回，一定要执行到write
func (ctx *Context) response() {
	code, resBytes := serializeRes(ctx.ResData)

	ctx.Writer.WriteHeader(code)

	_, err := ctx.Writer.Write(resBytes)
	if err != nil {
		mlog.Error("response failed", mlog.Field("error", err))
		return
	}
}

func serializeRes(obj any) (int, []byte) {
	code := http.StatusOK

	switch v := obj.(type) {
	case *Error:
		obj = &Response{Err: v.Error()}

		code = v.HttpCode
	default: // *api.resStruct(s)
		obj = &Response{
			IsSuccess: true,
			Data:      v,
		}
	}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		// 因为这里已经给resBytes定型了，返回错误也没啥能做的，就不返回了
		mlog.Error("serialize res to json failed", mlog.Field("error", err))
		return code, nil
	}

	return code, jsonBytes
}
