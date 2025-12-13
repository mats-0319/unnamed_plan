package main

import (
	"fmt"

	mconst "github.com/mats0319/unnamed_plan/server/internal/const"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

var serverNames = make(map[string]string) // uri - server name

func main() {
	mhttp.StartServer(newHandler())
}

func newHandler() *mhttp.Handler {
	h := &mhttp.Handler{}

	register("/api/login", mconst.ServerName_User, h)
	register("/api/user/list", mconst.ServerName_User, h)
	register("/api/user/create", mconst.ServerName_User, h)
	register("/api/user/modify", mconst.ServerName_User, h)
	register("/api/user/authenticate", mconst.ServerName_User, h)

	return h
}

func register(uri string, serverName string, h *mhttp.Handler) {
	serverNames[uri] = serverName
	h.AddHandler(uri, forward)
}

func forward(ctx *mhttp.Context) {
	uri := ctx.Request.RequestURI
	v, ok := serverNames[uri]

	newHost := ""
	switch v {
	case mconst.ServerName_User:
		newHost = "127.0.0.1:10320"
	}

	if !ok || len(newHost) < 1 {
		str := fmt.Sprintf("unknown req uri or server name: %s, %s", uri, v)
		mlog.Log("get server failed", str)
		ctx.ResData = str
		return
	}

	url := "http://" + newHost + uri
	err := ctx.Forward(url, ctx.Request.Body)
	if err != nil {
		mlog.Log("forward http req failed", mlog.Field("error", err))
		ctx.ResData = err
		return
	}
}
