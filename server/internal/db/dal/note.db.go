package dal

import (
	"context"
	"strings"

	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
	"gorm.io/gorm"
)

func GetNote(noteID string) (*model.Note, *Error) {
	res, err := Note.WithContext(context.TODO()).Where(Note.NoteID.Eq(noteID)).First()
	if err != nil {
		var e *Error
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			e = ErrNoteNotFound().WithCause(err)
		} else {
			e = ErrDBError().WithCause(err)
		}
		mlog.Error(e.String())
		return nil, e
	}

	return res, nil
}

func CreateNote(note *model.Note) *Error {
	if err := Note.WithContext(context.TODO()).Create(note); err != nil {
		var e *Error
		if strings.Contains(err.Error(), "violates unique constraint") {
			e = ErrNoteExist().WithCause(err)
		} else {
			e = ErrDBError().WithCause(err)
		}

		mlog.Error(e.String())
		return e
	}

	return nil
}

func UpdateNote(note *model.Note) *Error {
	if err := Note.WithContext(context.TODO()).Where(Note.ID.Eq(note.ID)).Save(note); err != nil {
		e := ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	return nil
}

func DeleteNote(noteID string) *Error {
	_, err := Note.WithContext(context.TODO()).Where(Note.NoteID.Eq(noteID)).Delete()
	if err != nil {
		e := ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	return nil
}

func ListNote(pageSize int, pageNum int, writer string) (int64, []*model.Note, *Error) {
	sql := Note.WithContext(context.TODO())
	if len(writer) > 0 {
		sql = sql.Where(Note.Writer.Eq(writer))
	}

	amount, err := sql.Count()
	if err != nil {
		e := ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return 0, nil, e
	}

	notes, err := sql.Order(Note.UpdatedAt.Desc()).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find()
	if err != nil {
		e := ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return 0, nil, e
	}

	return amount, notes, nil
}
