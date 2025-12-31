package main

import (
	"io"
	"strconv"

	"github.com/mats0319/unnamed_plan/server/gateway/middleware"
	"github.com/mats0319/unnamed_plan/server/internal/const"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

var serverNames = make(map[string]string) // uri - server name
var h = &mhttp.Handler{}

func main() {
	initHandler()
	mhttp.StartServer(h)
}

func initHandler() {
	register("/api"+api.URI_Login, mconst.ServerName_User)
	register("/api"+api.URI_Register, mconst.ServerName_User)
	register("/api"+api.URI_ListUser, mconst.ServerName_User, middleware.VerifyToken)
	register("/api"+api.URI_ModifyUser, mconst.ServerName_User, middleware.VerifyToken)
	register("/api"+api.URI_Authenticate, mconst.ServerName_User, middleware.VerifyToken)

	register("/api"+api.URI_CreateNote, mconst.ServerName_Note, middleware.VerifyToken)
	register("/api"+api.URI_ListNote, mconst.ServerName_Note)
	register("/api"+api.URI_ModifyNote, mconst.ServerName_Note, middleware.VerifyToken)
	register("/api"+api.URI_DeleteNote, mconst.ServerName_Note, middleware.VerifyToken)
}

func register(uri string, serverName string, middlewares ...func(ctx *mhttp.Context) error) {
	serverNames[uri] = serverName
	h.AddHandler(uri, forward, middlewares...)
}

func forward(ctx *mhttp.Context) {
	url, err := getURL(ctx)
	if err != nil {
		ctx.ResData = err
		return
	}

	res, err := ctx.Forward(url)
	if err != nil {
		ctx.ResData = err
		return
	}

	setLoginToken(ctx)

	// 这里我们希望统一使用ctx.response设置响应头和返回值，所以不使用io.copy直接复制res.body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		e := NewError(ET_ServerInternalError).WithCause(err)
		mlog.Log(e.String())
		ctx.ResData = e
		return
	}

	ctx.ResData = bodyBytes
}

func getURL(ctx *mhttp.Context) (string, error) {
	uri := ctx.Request.RequestURI
	v, ok := serverNames[uri]

	newHost := ""
	switch v {
	case mconst.ServerName_User:
		newHost = "127.0.0.1:10320"
	case mconst.ServerName_Note:
		newHost = "127.0.0.1:10321"
	}

	if !ok || len(newHost) < 1 {
		err := NewError(ET_ServerInternalError).WithParam("uri", uri).WithParam("server name", v)
		mlog.Log(err.String())
		return "", err
	}

	return "http://" + newHost + uri, nil
}

func setLoginToken(ctx *mhttp.Context) {
	userIDStr := ctx.ResHeaders[mconst.HttpHeader_UserID]
	userID, _ := strconv.Atoi(userIDStr)
	token := ctx.ResHeaders[mconst.HttpHeader_AccessToken]

	if userID > 0 && len(token) > 0 {
		middleware.SetToken(uint(userID), token)
	}
}
