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
	res := httpInvoke(api.URI_DeleteNote, `{"note_id":"not exist"}`, accessToken_User)
	if res.IsSuccess || res.Err != utils.ErrNoteNotFound().Error() {
		return unknownError
	}

	return ""
}

func deleteNoteCase_NotWriter() string {
	res := httpInvoke(api.URI_DeleteNote, fmt.Sprintf(`{"note_id":"%s"}`, noteID), accessToken_User)
	if res.IsSuccess || res.Err != utils.ErrNeedOwner().Error() {
		return unknownError
	}

	return ""
}

func deleteNoteCase_Success() string {
	res := httpInvoke(api.URI_DeleteNote, fmt.Sprintf(`{"note_id":"%s"}`, noteID), accessToken_Admin)
	if !res.IsSuccess {
		return res.Err
	}

	count, _, err := dal.ListNote(10, 1, "")
	if err != nil || count != 0 {
		return unknownError
	}

	return ""
}
