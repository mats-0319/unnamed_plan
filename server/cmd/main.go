package main

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/cmd/handlers"
	mconfig "github.com/mats0319/unnamed_plan/server/internal/config"
	mdb "github.com/mats0319/unnamed_plan/server/internal/db"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	"github.com/mats0319/unnamed_plan/server/internal/http/middleware"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

func main() {
	mconfig.Initialize()
	mlog.Initialize()
	mdb.Initialize()

	mhttp.StartServer(newHandler())
}

func newHandler() *mhttp.Handler {
	h := &mhttp.Handler{}

	uriPrefix := "/api"
	// user
	h.AddHandler(uriPrefix+api.URI_Register, handlers.Register)
	h.AddHandler(uriPrefix+api.URI_Login, handlers.Login)
	h.AddHandler(uriPrefix+api.URI_ListUser, handlers.ListUser, middleware.VerifyToken)
	h.AddHandler(uriPrefix+api.URI_ModifyUser, handlers.ModifyUser, middleware.VerifyToken)

	// note
	h.AddHandler(uriPrefix+api.URI_CreateNote, handlers.CreateNote, middleware.VerifyToken)
	h.AddHandler(uriPrefix+api.URI_ListNote, handlers.ListNote)
	h.AddHandler(uriPrefix+api.URI_ModifyNote, handlers.ModifyNote, middleware.VerifyToken)
	h.AddHandler(uriPrefix+api.URI_DeleteNote, handlers.DeleteNote, middleware.VerifyToken)

	return h
}
