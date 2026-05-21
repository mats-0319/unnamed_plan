package mhttp

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

type Handler struct {
	Handlers map[string]*HandlerItem // uri - handler func and middleware(s)
	uri      []string                // 按照添加顺序保存uri，用于启动时展示
}

type HandlerItem struct {
	Func        func(ctx *Context)
	Middlewares []func(ctx *Context) *utils.Error
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "*")
	writer.Header().Set("Access-Control-Expose-Headers", "*")

	if request.Method != http.MethodPost { // return after set header
		return // res body is empty
	}

	handlerItemIns, ok := h.Handlers[request.RequestURI]
	if !ok {
		// 实际场景中，会触发这一段代码的，更多的是恶意向服务器批量发送请求的，所以就直接返回了。
		// 想看访问细节可以去nginx访问日志。
		// 不想在nginx上配置接口白名单，那样太死板了。找找有什么好办法。
		return
	}

	mlog.Info(fmt.Sprintf("> Receive Request: %s .", request.URL.String()))
	startTime := time.Now()
	defer func() {
		mlog.Info(fmt.Sprintf("> Process Request: %s , in %d ms", request.URL.String(), time.Since(startTime).Milliseconds()))

		if err := recover(); err != nil {
			mlog.Error("recover panic", slog.Any("", err))
		}
	}()

	ctx := NewContext(writer, request)
	defer ctx.response()

	// middlewares
	for i := range handlerItemIns.Middlewares {
		if err := handlerItemIns.Middlewares[i](ctx); err != nil { // log in middleware
			ctx.ResData = err
			return
		}
	}

	// business handler
	handlerItemIns.Func(ctx)
}

func (h *Handler) AddHandler(uri string, handlerFunc func(ctx *Context), middlewares ...func(ctx *Context) *utils.Error) {
	if h.Handlers == nil {
		h.Handlers = make(map[string]*HandlerItem)
	}

	h.Handlers[uri] = &HandlerItem{
		Func:        handlerFunc,
		Middlewares: middlewares,
	}

	h.uri = append(h.uri, uri)
}

func (h *Handler) DisplayRegisteredURI() {
	mlog.Info(fmt.Sprintf("> HTTP Handler Count: %d.", len(h.Handlers)))

	for _, v := range h.uri {
		mlog.Info("- HTTP Handler Registered: " + v)
	}
}
