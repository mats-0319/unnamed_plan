package api

import (
	"fmt"

	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

func ModifyNote() {
	testCase("note not exist", modifyNoteCase_NoteNotExist)
	testCase("no changes", modifyNoteCase_NoChanges)
	testCase("not writer", modifyNoteCase_NotWriter)
	testCase("success", modifyNoteCase_Success)
}

func modifyNoteCase_NoteNotExist() string {
	res := httpInvoke(api.URI_ModifyNote, `{"note_id":"not exist","is_anonymous":false,"title":"","content":"1"}`, accessToken_User)
	if res.IsSuccess || res.Err != utils.ErrNoteNotFound().Error() {
		return unknownError
	}

	return ""
}

func modifyNoteCase_NoChanges() string {
	res := httpInvoke(api.URI_ModifyNote, fmt.Sprintf(`{"note_id":"%s","is_anonymous":false,"title":"123","content":"456"}`, noteID), accessToken_Admin)
	if res.IsSuccess || res.Err != utils.ErrNoChanges().Error() {
		return unknownError
	}

	return ""
}

func modifyNoteCase_NotWriter() string {
	res := httpInvoke(api.URI_ModifyNote, fmt.Sprintf(`{"note_id":"%s","is_anonymous":false,"title":"","content":"1"}`, noteID), accessToken_User)
	if res.IsSuccess || res.Err != utils.ErrPermissionDenied().Error() {
		return unknownError
	}

	return ""
}

func modifyNoteCase_Success() string {
	res := httpInvoke(api.URI_ModifyNote, fmt.Sprintf(`{"note_id":"%s","is_anonymous":true,"title":"123123","content":"456456"}`, noteID), accessToken_Admin)
	if !res.IsSuccess {
		return res.Err
	}

	data, err := dal.GetNote(noteID)
	if data == nil || !data.IsAnonymous || data.Title != "123123" || data.Content != "456456" || err != nil {
		return unknownError
	}

	return ""
}
