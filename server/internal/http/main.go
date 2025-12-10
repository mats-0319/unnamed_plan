package mhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	mconf "github.com/mats0319/unnamed_plan/server/internal/config"
	mconst "github.com/mats0319/unnamed_plan/server/internal/const"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

type handlerFunc func(ctx *Context)

type Handler struct {
	handlers map[string]handlerFunc // pattern - handler func
}

type Config struct {
	Port string `json:"port"`
}

// StartServer is blocked
func StartServer(handler *Handler) {
	configIns := getConfig()

	addr := "127.0.0.1:" + configIns.Port
	mlog.Log("> Listening at: " + addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		mlog.Log("handlers listen and serve failed", mlog.Field("error", err))
	}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "*")
	writer.Header().Set("Access-Control-Expose-Headers", "*")

	mlog.Log(fmt.Sprintf("| %s | %s |", time.Now().String(), request.URL.String()))
	// todo: log res

	ctx := NewContext(writer, request)

	if ctx.Request.Method != http.MethodPost { // only accept 'post' req
		ctx.response()
		return
	}

	handlerFuncIns, ok := h.handlers[request.RequestURI]
	if !ok {
		ctx.ResData = mlog.Field("unsupported request uri", request.RequestURI)
		mlog.Log(ctx.ResData.(string))
		ctx.response()
		return
	}

	handlerFuncIns(ctx)

	ctx.response()
}

func (h *Handler) AddHandler(uri string, handler handlerFunc) {
	if h.handlers == nil {
		h.handlers = make(map[string]handlerFunc)
	}

	h.handlers[uri] = handler
}

func getConfig() *Config {
	jsonBytes := mconf.GetConfigItem(mconst.UID_Http)

	res := &Config{}
	err := json.Unmarshal(jsonBytes, res)
	if err != nil {
		mlog.Log("get gateway config failed", mlog.Field("error", err))
	}

	return res
}
