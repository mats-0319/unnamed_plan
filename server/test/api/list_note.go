package api

import (
	"log"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
)

func ListNote() {
	TestApi("List Note")

	TestCase("success")
	res := HttpInvoke(api.URI_ListNote, `{"page":{"size":10,"num":1},"user_id":0}`)
	log.Println(res)

	TestApiEnd()
}
