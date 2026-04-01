package api

import (
	"context"

	api "github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

var noteID string

func CreateNote() {
	testCase("success", createNoteCase_Success)
	testCase("duplicate create", createNoteCase_Duplicate)
}

func createNoteCase_Success() string {
	res := httpInvoke(api.URI_CreateNote, `{"is_anonymous":false,"title":"123","content":"456"}`, accessToken_Admin)
	if !res.IsSuccess {
		return res.Err
	}

	// record 'note id' for later api(s)
	note, _ := dal.Note.WithContext(context.TODO()).First()
	if note == nil {
		return unknownError
	}
	noteID = note.NoteID

	data, err := dal.GetNote(noteID)
	if err != nil || data == nil || data.Title != "123" || data.Content != "456" {
		return unknownError
	}

	return ""
}

func createNoteCase_Duplicate() string {
	res := httpInvoke(api.URI_CreateNote, `{"is_anonymous":false,"title":"123","content":"456"}`, accessToken_Admin)
	if res.IsSuccess || res.Err != utils.ErrNoteExist().Error() {
		return unknownError
	}

	return ""
}
