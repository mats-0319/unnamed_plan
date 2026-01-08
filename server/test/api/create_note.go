package api

import (
	"log"

	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
)

func CreateNote() {
	TestApi("Create Note")
	HttpInvoke(api.URI_Login, `{"user_name":"user","password":"8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92","totp_code":""}`)

	// 通常不会重复创建，因为只有一个note id参数要求唯一，但这个参数还是服务端随机生成的

	TestCase("success")
	res := HttpInvoke(api.URI_CreateNote, `{"is_anonymous":false,"title":"123","content":"456"}`)
	log.Println(res)

	TestApiEnd()
}
