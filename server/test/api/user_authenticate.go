package api

import (
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
)

func UserAuthenticate() {
	TestApi("Authenticate User")

	// success
	HttpInvoke(api.URI_Authenticate, ``)

	TestApiEnd()
}
