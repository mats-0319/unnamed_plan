package main

import (
	"github.com/mats0319/unnamed_plan/server/1_user/handlers"
	_ "github.com/mats0319/unnamed_plan/server/internal/db"
	mhttp "github.com/mats0319/unnamed_plan/server/internal/http"
)

func main() {
	mhttp.StartServer(newHandler())
}

func newHandler() *mhttp.Handler {
	h := &mhttp.Handler{}

	h.AddHandler("/api/login", handlers.Login)

	return h
}
