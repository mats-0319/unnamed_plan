package main

import (
	"github.com/mats0319/unnamed_plan/server/1_user/handlers"
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
	h.AddHandler("/api"+api.URI_Register, handlers.Register)
	h.AddHandler("/api"+api.URI_ListUser, handlers.ListUser)
	h.AddHandler("/api"+api.URI_ModifyUser, handlers.ModifyUser)
	h.AddHandler("/api"+api.URI_Authenticate, handlers.Authenticate)

	return h
}
