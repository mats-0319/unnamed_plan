package api

import (
	"fmt"

	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func DeleteNote() {
	testCase("note not exist", deleteNoteCase_NoteNotExist)
	testCase("not writer", deleteNoteCase_NotWriter)
	testCase("success", deleteNoteCase_Success)
}

func deleteNoteCase_NoteNotExist() string {
	res := httpInvoke(api.URI_DeleteNote, `{"note_id":"not exist"}`)
	if res.IsSuccess || res.Err != utils.ErrNoteNotFound().Error() {
		return unknownError
	}

	return ""
}

func deleteNoteCase_NotWriter() string {
	loginCase_Success(false)()

	res := httpInvoke(api.URI_DeleteNote, fmt.Sprintf(`{"note_id":"%s"}`, noteID))
	if res.IsSuccess || res.Err != utils.ErrNeedOwner().Error() {
		return unknownError
	}

	return ""
}

func deleteNoteCase_Success() string {
	loginCase_Success(true)()

	res := httpInvoke(api.URI_DeleteNote, fmt.Sprintf(`{"note_id":"%s"}`, noteID))
	if !res.IsSuccess {
		return res.Err
	}

	count, _, err := dal.ListNote(api.Pagination{Size: 10, Num: 1}, "")
	if count != 0 || err != nil {
		return unknownError
	}

	return ""
}
