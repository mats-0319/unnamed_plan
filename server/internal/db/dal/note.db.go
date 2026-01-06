package dal

import (
	"context"
	"strings"

	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

func GetNote(id uint) (*model.Note, *Error) {
	qn := Q.Note
	sql := qn.WithContext(context.TODO()).Where(qn.ID.Eq(id))

	res, err := sql.First()
	if err != nil {
		var e *Error
		if strings.Contains(err.Error(), "record not found") {
			e = NewError(ET_ParamsError, ED_NoteNotFound).WithCause(err)
		} else {
			e = NewError(ET_ServerInternalError).WithCause(err)
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
			e = NewError(ET_ParamsError, ED_NoteExist).WithCause(err)
		} else {
			e = NewError(ET_ServerInternalError).WithCause(err)
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
		e := NewError(ET_ServerInternalError).WithCause(err)
		mlog.Log(e.String())
		return e
	}

	return nil
}

func DeleteNote(id uint) *Error {
	qn := Q.Note
	_, err := qn.WithContext(context.TODO()).Where(qn.ID.Eq(id)).Delete()
	if err != nil {
		e := NewError(ET_ServerInternalError).WithCause(err)
		mlog.Log(e.String())
		return e
	}

	return nil
}

func ListNote(page api.Pagination, writerID ...uint) (int64, []*model.Note, *Error) {
	qn := Q.Note
	sql := qn.WithContext(context.TODO())
	if len(writerID) > 0 {
		sql = sql.Where(qn.WriterID.Eq(writerID[0]))
	}

	amount, err := sql.Count()
	if err != nil {
		e := NewError(ET_ServerInternalError).WithCause(err)
		mlog.Log(e.String())
		return 0, nil, e
	}

	res, err := sql.Order(qn.UpdatedAt.Desc()).Limit(page.Size).Offset((page.Num - 1) * page.Size).Find()
	if err != nil {
		e := NewError(ET_ServerInternalError).WithCause(err)
		mlog.Log(e.String())
		return 0, nil, e
	}

	return amount, res, nil
}
