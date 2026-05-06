package dal

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
	"gorm.io/gorm"
)

func GetNote(noteID string) (res *model.Note, e *Error) {
	res, err := Note.WithContext(context.TODO()).Where(Note.NoteID.Eq(noteID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e = ErrNoteNotFound().WithCause(err)
		} else {
			e = ErrDBError().WithCause(err)
		}

		mlog.Error(e.String())
		return
	}

	return
}

func CreateNote(note *model.Note) (e *Error) {
	if err := Note.WithContext(context.TODO()).Create(note); err != nil {
		if pge, ok := errors.AsType[*pgconn.PgError](err); ok && pge.Code == "23505" {
			e = ErrNoteExist().WithCause(err)
		} else {
			e = ErrDBError().WithCause(err)
		}

		mlog.Error(e.String())
		return
	}

	return
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
	if _, err := Note.WithContext(context.TODO()).Where(Note.NoteID.Eq(noteID)).Delete(); err != nil {
		e := ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	return nil
}

func ListNote(pageSize int, pageNum int, writer string) (count int64, records []*model.Note, e *Error) {
	sql := Note.WithContext(context.TODO())
	if len(writer) > 0 {
		sql = sql.Where(Note.Writer.Eq(writer))
	}

	count, err := sql.Count()
	if err != nil {
		e = ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return
	}

	records, err = sql.Order(Note.UpdatedAt.Desc()).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find()
	if err != nil {
		e = ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return
	}

	return
}
