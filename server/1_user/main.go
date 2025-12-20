package main

import (
	"github.com/mats0319/unnamed_plan/server/1_user/handlers"
	"github.com/mats0319/unnamed_plan/server/1_user/middleware"
	_ "github.com/mats0319/unnamed_plan/server/internal/db"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
)

func main() {
	mhttp.StartServer(newHandler())
}

func newHandler() *mhttp.Handler {
	h := &mhttp.Handler{}

	h.AddHandler("/api"+api.URI_Login, handlers.Login)
	h.AddHandler("/api"+api.URI_CreateUser, handlers.CreateUser)
	h.AddHandler("/api"+api.URI_ListUser, handlers.ListUser, middleware.VerifyToken)
	h.AddHandler("/api"+api.URI_ModifyUser, handlers.ModifyUser, middleware.VerifyToken)
	h.AddHandler("/api"+api.URI_Authenticate, handlers.Authenticate, middleware.VerifyToken)

	return h
}
