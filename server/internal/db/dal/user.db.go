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

func GetUser(userName string) (res *model.User, e *Error) {
	res, err := User.WithContext(context.TODO()).Where(User.UserName.Eq(userName)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e = ErrUserNotFound().WithCause(err)
		} else {
			e = ErrDBError().WithCause(err)
		}

		mlog.Error(e.String())
		return
	}

	return
}

func CreateUser(user *model.User) (e *Error) {
	if err := User.WithContext(context.TODO()).Create(user); err != nil {
		if pge, ok := errors.AsType[*pgconn.PgError](err); ok && pge.Code == "23505" {
			e = ErrUserExist().WithCause(err)
		} else {
			e = ErrDBError().WithCause(err)
		}

		mlog.Error(e.String())
		return
	}

	return
}

func UpdateUser(user *model.User) *Error {
	if err := User.WithContext(context.TODO()).Where(User.ID.Eq(user.ID)).Save(user); err != nil {
		e := ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	return nil
}

func ListUser(pageSize int, pageNum int) (count int64, records []*model.User, e *Error) {
	sql := User.WithContext(context.TODO())

	count, err := sql.Count()
	if err != nil {
		e = ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return
	}

	records, err = sql.Order(User.UpdatedAt.Desc()).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find()
	if err != nil {
		e = ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return
	}

	return
}
