package main

import (
	"io"
	"strconv"

	"github.com/mats0319/unnamed_plan/server/gateway/middleware"
	. "github.com/mats0319/unnamed_plan/server/internal/const"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

var serverNames = make(map[string]string) // uri - server name
var h = &mhttp.Handler{}

func main() {
	initHandler()
	mhttp.StartServer(h)
}

func initHandler() {
	register("/api"+api.URI_Login, ServerName_User)
	register("/api"+api.URI_Register, ServerName_User)
	register("/api"+api.URI_ListUser, ServerName_User, middleware.VerifyToken)
	register("/api"+api.URI_ModifyUser, ServerName_User, middleware.VerifyToken)
	register("/api"+api.URI_Authenticate, ServerName_User, middleware.VerifyToken)

	register("/api"+api.URI_CreateNote, ServerName_Note, middleware.VerifyToken)
	register("/api"+api.URI_ListNote, ServerName_Note)
	register("/api"+api.URI_ModifyNote, ServerName_Note, middleware.VerifyToken)
	register("/api"+api.URI_DeleteNote, ServerName_Note, middleware.VerifyToken)
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

	reader, err := ctx.Forward(url)
	if err != nil {
		ctx.ResData = err
		return
	}

	setLoginToken(ctx)

	// 这里我们希望统一使用ctx.response设置响应头和返回值，所以不使用io.copy直接复制res.body
	bodyBytes, err := io.ReadAll(reader)
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
	case ServerName_User:
		newHost = "127.0.0.1:10320"
	case ServerName_Note:
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
	userIDStr := ctx.ResHeaders[HttpHeader_UserID]
	userID, _ := strconv.Atoi(userIDStr)
	token := ctx.ResHeaders[HttpHeader_AccessToken]

	if userID > 0 && len(token) > 0 {
		middleware.SetToken(uint(userID), token)
	}
}
