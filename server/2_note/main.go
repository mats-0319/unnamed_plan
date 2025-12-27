package main

import (
	"github.com/mats0319/unnamed_plan/server/2_note/handlers"
	_ "github.com/mats0319/unnamed_plan/server/internal/db"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
)

func main() {
	mhttp.StartServer(newHandler())
}

func newHandler() *mhttp.Handler {
	h := &mhttp.Handler{}

	h.AddHandler("/api"+api.URI_CreateNote, handlers.CreateNote)
	h.AddHandler("/api"+api.URI_ListNote, handlers.ListNote)
	h.AddHandler("/api"+api.URI_ModifyNote, handlers.ModifyNote)
	h.AddHandler("/api"+api.URI_DeleteNote, handlers.DeleteNote)

	return h
}
