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
	loginCase_Success(true)()

	res := httpInvoke(api.URI_CreateNote, `{"is_anonymous":false,"title":"123","content":"456"}`)
	if !res.IsSuccess {
		return res.Err
	}

	// get note id from db
	note, _ := dal.Q.Note.WithContext(context.TODO()).First()
	if note == nil {
		return unknownError
	}
	noteID = note.NoteID

	data, err := dal.GetNote(noteID)
	if data == nil || data.Title != "123" || data.Content != "456" || err != nil {
		return unknownError
	}

	return ""
}

func createNoteCase_Duplicate() string {
	res := httpInvoke(api.URI_CreateNote, `{"is_anonymous":false,"title":"123","content":"456"}`)
	if res.IsSuccess || res.Err != utils.ErrNoteExist().Error() {
		return unknownError
	}

	return ""
}
