package api

import (
	"log"

	api2 "github.com/mats0319/unnamed_plan/server/cmd/api/go"
)

func DeleteNote() {
	TestApi("Delete Note")

	TestCase("note not exist")
	HttpInvoke(api2.URI_DeleteNote, `{"id":0}`)

	TestCase("not writer")
	HttpInvoke(api2.URI_Login, `{"user_name":"mats0319","password":"8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92","totp_code":""}`)
	HttpInvoke(api2.URI_DeleteNote, `{"id":1003}`)

	TestCase("success")
	HttpInvoke(api2.URI_Login, `{"user_name":"user","password":"8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92","totp_code":""}`)
	res := HttpInvoke(api2.URI_DeleteNote, `{"id":1003}`)
	log.Println(res)

	TestApiEnd()
}
