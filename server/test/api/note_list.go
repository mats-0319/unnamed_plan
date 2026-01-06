package api

import (
	"log"

	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
)

func NoteList() {
	TestApi("List Note")

	TestCase("success")
	res := HttpInvoke(api.URI_ListNote, `{"page":{"size":10,"num":1},"list_my_flag":false}`)
	log.Println(res)

	TestApiEnd()
}
