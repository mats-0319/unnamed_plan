package api

import (
	"log"

	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
)

func DeleteNote() {
	TestApi("Delete Note")

	TestCase("note not exist")
	HttpInvoke(api.URI_DeleteNote, `{"id":0}`)

	TestCase("not writer")
	HttpInvoke(api.URI_Login, `{"user_name":"mats0319","password":"8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92","totp_code":""}`)
	HttpInvoke(api.URI_DeleteNote, `{"id":1001}`)

	TestCase("success")
	HttpInvoke(api.URI_Login, `{"user_name":"user","password":"8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92","totp_code":""}`)
	res := HttpInvoke(api.URI_DeleteNote, `{"id":1001}`)
	log.Println(res)

	TestApiEnd()
}
