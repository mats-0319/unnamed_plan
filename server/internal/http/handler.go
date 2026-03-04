package mhttp

import (
	"fmt"
	"net/http"
	"time"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

type Handler struct {
	handlers map[string]*HandlerItem // uri - handler func
}

type HandlerItem struct {
	Func        func(ctx *Context)
	Middlewares []func(ctx *Context) *Error
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "*")
	writer.Header().Set("Access-Control-Expose-Headers", "*")

	if request.Method != http.MethodPost { // return after set header
		return // res body is empty
	}

	mlog.Info(fmt.Sprintf("| %s | %s |", request.URL.String(), time.Now().String()))

	ctx := NewContext(writer, request)
	defer ctx.response()

	handlerItemIns, ok := h.handlers[request.RequestURI]
	if !ok {
		e := ErrUnknownUri().WithParam("uri", request.RequestURI)
		mlog.Error(e.String())
		ctx.ResData = e
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

	// business handler
	handlerItemIns.Func(ctx)
}

func (h *Handler) AddHandler(uri string, handlerFunc func(ctx *Context), middlewares ...func(ctx *Context) *Error) {
	if h.handlers == nil {
		h.handlers = make(map[string]*HandlerItem)
	}

	h.handlers[uri] = &HandlerItem{
		Func:        handlerFunc,
		Middlewares: middlewares,
	}
}

func (h *Handler) displayRegisteredUri() {
	for k := range h.handlers {
		mlog.Info("- Http Handler Registered: " + k)
	}
}
