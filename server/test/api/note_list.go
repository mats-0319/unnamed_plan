package api

import (
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
)

func NoteList() {
	TestApi("List Note")

	// success
	HttpInvoke(api.URI_ListNote, `{"page":{"size":10,"num":1},"list_my_flag":false}`)

	TestApiEnd()
}
