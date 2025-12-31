package mhttp

import (
	"fmt"
	"net/http"
	"time"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

type Handler struct {
	config   *config
	handlers map[string]*HandlerItem // uri - handler func
}

type HandlerItem struct {
	Func        func(ctx *Context)
	Middlewares []func(ctx *Context) error
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := NewContext(writer, request)
	defer ctx.response()

	origin := request.Header.Get("Origin")
	if !h.isValidOrigin(origin) {
		ctx.ResData = NewError(ET_ServerInternalError).WithParam("origin", origin)
		return
	}

	writer.Header().Set("Access-Control-Allow-Origin", origin)
	writer.Header().Set("Access-Control-Allow-Headers", "*")
	writer.Header().Set("Access-Control-Expose-Headers", "*")

	if request.Method != http.MethodPost { // return after set header
		return
	}

	mlog.Log(fmt.Sprintf("| %s | %s |", request.URL.String(), time.Now().String()))

	handlerItemIns, ok := h.handlers[request.RequestURI]
	if !ok {
		err := NewError(ET_ServerInternalError).WithParam("uri", request.RequestURI)
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

func (h *Handler) isValidOrigin(origin string) bool {
	containsFlag := false
	for _, v := range h.config.AllowedOrigins {
		if origin == v || v == "*" {
			containsFlag = true
			break
		}
	}

	return containsFlag
}
