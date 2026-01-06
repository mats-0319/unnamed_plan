package api

import (
	"log"

	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
)

func UserAuthenticate() {
	TestApi("Authenticate User")

	TestCase("success")
	res := HttpInvoke(api.URI_Authenticate, ``)
	log.Println(res)

	TestApiEnd()
}
