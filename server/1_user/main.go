package main

import (
	"github.com/mats0319/unnamed_plan/server/1_user/handlers"
	"github.com/mats0319/unnamed_plan/server/1_user/middleware"
	_ "github.com/mats0319/unnamed_plan/server/internal/db"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
)

func main() {
	mhttp.StartServer(newHandler())
}

func newHandler() *mhttp.Handler {
	h := &mhttp.Handler{}

	h.AddHandler("/api/login", handlers.Login)
	h.AddHandler("/api/user/list", handlers.ListUser, middleware.VerifyToken)
	h.AddHandler("/api/user/create", handlers.CreateUser, middleware.VerifyToken)
	h.AddHandler("/api/user/modify", handlers.ModifyUser, middleware.VerifyToken)
	h.AddHandler("/api/user/authenticate", handlers.Authenticate, middleware.VerifyToken)

	return h
}
