package api

import (
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
)

func UserLogin() {
	TestApi("Login")

	TestCase("user not exist")
	HttpInvoke(api.URI_Login, `{"user_name":"not exist","password":"","totp_code":""}`)

	TestCase("wrong pwd")
	HttpInvoke(api.URI_Login, `{"user_name":"mats0319","password":"wrong pwd","totp_code":""}`)

	TestCase("wrong totp code")
	HttpInvoke(api.URI_Login, `{
"user_name":"admin",
"password":"8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92",
"totp_code":"000000"}`)

	// success
	HttpInvoke(api.URI_Login, `{"user_name":"user","password":"8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92","totp_code":""}`)

	TestApiEnd()
}
