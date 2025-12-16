package main

import (
	. "github.com/mats0319/unnamed_plan/server/internal/const"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

var serverNames = make(map[string]string) // uri - server name

func main() {
	mhttp.StartServer(newHandler())
}

func newHandler() *mhttp.Handler {
	h := &mhttp.Handler{}

	register("/api/login", ServerName_User, h)
	register("/api/user/list", ServerName_User, h)
	register("/api/user/create", ServerName_User, h)
	register("/api/user/modify", ServerName_User, h)
	register("/api/user/authenticate", ServerName_User, h)

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
	case ServerName_User:
		newHost = "127.0.0.1:10320"
	}

	if !ok || len(newHost) < 1 {
		err := NewError(ET_ServerInternalError, ED_UnknownURIOrServerName).
			WithParam("uri", uri).WithParam("server name", v)
		mlog.Log(err.String())
		ctx.ResData = err
		return
	}

	url := "http://" + newHost + uri
	ctx.Forward(url, ctx.Request.Body)
}
