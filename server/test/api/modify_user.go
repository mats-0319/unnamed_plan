package api

import (
	"log"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
)

func ModifyUser() {
	TestApi("Modify User")
	HttpInvoke(api.URI_Login, `{"user_name":"user","password":"8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92","totp_code":""}`)

	TestCase("no changes")
	HttpInvoke(api.URI_ModifyUser, `{"nickname":"","password":"","modify_tk_flag":false,"totp_key":""}`)

	// operator not exist(?): need mock

	TestCase("same pwd")
	HttpInvoke(api.URI_ModifyUser, `{"nickname":"",
"password":"8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92",
"modify_tk_flag":false,"totp_key":""}`)

	TestCase("invalid totp key")
	HttpInvoke(api.URI_ModifyUser, `{"nickname":"","password":"","modify_tk_flag":true,"totp_key":"123"}`)

	TestCase("success")
	res := HttpInvoke(api.URI_ModifyUser, `{"nickname":"123","password":"","modify_tk_flag":false,"totp_key":""}`)
	log.Println(res)

	TestApiEnd()
}
