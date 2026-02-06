package api

import (
	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
)

func CreateNote() {
	testCase("success", createNoteCase_Success)
}

func createNoteCase_Success() string {
	loginCase_Success(true)()

	res := httpInvoke(api.URI_CreateNote, `{"is_anonymous":false,"title":"123","content":"456"}`)
	if !res.IsSuccess {
		return res.Err
	}

	data, err := dal.GetNote(1001)
	if data == nil || data.Title != "123" || data.Content != "456" || err != nil {
		return unknownError
	}

	return ""
}
