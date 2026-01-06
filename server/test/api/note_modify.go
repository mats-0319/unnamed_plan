package api

import api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"

func NoteModify() {
	TestApi("Modify Note")

	TestCase("note not exist")
	HttpInvoke(api.URI_ModifyNote, `{"id":0,"is_anonymous":false,"title":"","content":""}`)

	TestCase("not writer")
	HttpInvoke(api.URI_Login, `{"user_name":"mats0319","password":"8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92","totp_code":""}`)
	HttpInvoke(api.URI_ModifyNote, `{"id":1001,"is_anonymous":false,"title":"","content":""}`) // 这里借用create接口创建的note

	TestCase("no changes")
	HttpInvoke(api.URI_Login, `{"user_name":"user","password":"8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92","totp_code":""}`)
	HttpInvoke(api.URI_ModifyNote, `{"id":1001,"is_anonymous":false,"title":"123","content":"456"}`)

	// success
	HttpInvoke(api.URI_ModifyNote, `{"id":1001,"is_anonymous":true,"title":"123","content":"456"}`)

	TestApiEnd()
}
