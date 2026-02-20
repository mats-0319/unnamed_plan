package api

import (
	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
)

func ListNote() {
	testCase("success", listNoteCase_Success)
}

func listNoteCase_Success() string {
	res := httpInvoke(api.URI_ListNote, `{"page":{"size":10,"num":1},"user_id":0}`)
	if !res.IsSuccess {
		return res.Err
	}

	count, data, err := dal.ListNote(api.Pagination{Size: 10, Num: 1}, "")
	if count != 1 || len(data) != 1 || err != nil {
		return unknownError
	}

	// test list notes with given writer
	count, data, err = dal.ListNote(api.Pagination{Size: 10, Num: 1}, "not exist")
	if count != 0 || len(data) != 0 || err != nil {
		return unknownError
	}

	return ""
}
