package dal

import (
	"context"
	"strings"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/cmd/model"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

func GetNote(noteID string) (*model.Note, *Error) {
	qn := Q.Note
	sql := qn.WithContext(context.TODO()).Where(qn.NoteID.Eq(noteID))

	res, err := sql.First()
	if err != nil {
		var e *Error
		if strings.Contains(err.Error(), "record not found") {
			e = ErrNoteNotFound().WithCause(err)
		} else {
			e = ErrDBError().WithCause(err)
		}
		mlog.Log(e.String())
		return nil, e
	}

	return res, nil
}

func CreateNote(note *model.Note) *Error {
	err := Q.Note.WithContext(context.TODO()).Create(note)
	if err != nil {
		var e *Error
		if strings.Contains(err.Error(), "violates unique constraint") {
			e = ErrNoteExist().WithCause(err)
		} else {
			e = ErrDBError().WithCause(err)
		}

		mlog.Log(e.String())
		return e
	}

	return nil
}

func UpdateNote(note *model.Note) *Error {
	qn := Q.Note
	err := qn.WithContext(context.TODO()).Where(qn.ID.Eq(note.ID)).Save(note)
	if err != nil {
		e := ErrDBError().WithCause(err)
		mlog.Log(e.String())
		return e
	}

	return nil
}

func DeleteNote(noteID string) *Error {
	qn := Q.Note
	_, err := qn.WithContext(context.TODO()).Where(qn.NoteID.Eq(noteID)).Delete()
	if err != nil {
		e := ErrDBError().WithCause(err)
		mlog.Log(e.String())
		return e
	}

	return nil
}

func ListNote(page api.Pagination, writer string) (int64, []*model.Note, *Error) {
	qn := Q.Note
	sql := qn.WithContext(context.TODO())
	if len(writer) > 0 {
		sql = sql.Where(qn.Writer.Eq(writer))
	}

	amount, err := sql.Count()
	if err != nil {
		e := ErrDBError().WithCause(err)
		mlog.Log(e.String())
		return 0, nil, e
	}

	res, err := sql.Order(qn.UpdatedAt.Desc()).Limit(page.Size).Offset((page.Num - 1) * page.Size).Find()
	if err != nil {
		e := ErrDBError().WithCause(err)
		mlog.Log(e.String())
		return 0, nil, e
	}

	return amount, res, nil
}
