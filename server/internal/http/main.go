package mhttp

import (
	"net/http"

	mconst "github.com/mats0319/unnamed_plan/server/internal/const"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

type handlerFunc func(ctx *Context) *Context

type HttpHandler struct {
	Port string `json:"port"`

	handlers map[string]handlerFunc // pattern - handler func
}

// StartServer is blocked
func StartServer(handler *HttpHandler) {
	mlog.Log("> Listening at: 127.0.0.1:" + handler.Port)

	if err := http.ListenAndServe("127.0.0.1:"+handler.Port, handler); err != nil {
		mlog.Log("http listen and serve failed", mlog.Field[error]("error", err))
	}
}

func (h *HttpHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// todo: log req

	handlerFuncIns, ok := h.handlers[request.RequestURI]
	if !ok {
		str := mlog.Field(mconst.HttpErrType_UnsupportedReqURI, request.RequestURI)
		mlog.Log(str)
		response(writer, NewResData(str))
		return
	}

	ctx := NewContext(writer, request)

	c := handlerFuncIns(ctx)

	ctx.ResData = c.ResData

	// todo: log res (with invoke chain)

	response(writer, ctx.ResData)
}

func (h *HttpHandler) AddHandler(uri string, handler handlerFunc) {
	h.handlers[uri] = handler
}
