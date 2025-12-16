package mhttp

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/mats0319/unnamed_plan/server/internal/const"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

type Handler struct {
	handlers map[string]*HandlerItem // uri - handler func
}

type HandlerItem struct {
	Func        func(ctx *Context)
	Middlewares []func(ctx *Context) error
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "*")
	writer.Header().Set("Access-Control-Expose-Headers", "*")

	mlog.Log(fmt.Sprintf("| %s | %s |", request.URL.String(), time.Now().String()))

	ctx := NewContext(writer, request)
	defer ctx.response()

	if ctx.Request.Method != http.MethodPost { // only accept 'post' req
		return
	}

	handlerItemIns, ok := h.handlers[request.RequestURI]
	if !ok {
		err := NewError(ET_ServerInternalError, ED_UnsupportedURI).WithParam("uri", request.RequestURI)
		mlog.Log(err.String())
		ctx.ResData = err
		return
	}

	// middlewares
	for i := range handlerItemIns.Middlewares {
		err := handlerItemIns.Middlewares[i](ctx)
		if err != nil { // log in middleware
			ctx.ResData = err
			return
		}
	}

	handlerItemIns.Func(ctx)
}

func (h *Handler) AddHandler(uri string, handlerFunc func(ctx *Context), middlewares ...func(ctx *Context) error) {
	if h.handlers == nil {
		h.handlers = make(map[string]*HandlerItem)
	}

	h.handlers[uri] = &HandlerItem{
		Func:        handlerFunc,
		Middlewares: middlewares,
	}
}

func (h *Handler) supportedUri() {
	for k := range h.handlers {
		mlog.Log(fmt.Sprintf("| Http Registered: %s", k))
	}
}
