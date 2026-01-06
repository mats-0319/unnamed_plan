package api

import (
	"fmt"
	"log"
	"time"

	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
)

func UserRegister() {
	TestApi("Register")

	TestCase("duplicate register")
	HttpInvoke(api.URI_Register, `{"user_name":"admin","password":""}`)

	TestCase("success")
	newUserName := fmt.Sprintf("new user %d", time.Now().UnixMilli())
	res := HttpInvoke(api.URI_Register, fmt.Sprintf(`{"user_name":"%s","password":"123"}`, newUserName))
	log.Println(res)

	TestApiEnd()
}
