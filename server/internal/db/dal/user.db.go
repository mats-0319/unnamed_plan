package dal

import (
	"context"
	"strings"

	"github.com/mats0319/unnamed_plan/server/cmd/api/go"
	"github.com/mats0319/unnamed_plan/server/cmd/model"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
)

// GetUser query user by 'id'/'username', according to value type
func GetUser[T uint | string](value T) (*model.User, *Error) {
	qu := Q.User
	sql := qu.WithContext(context.TODO())

	switch v := any(value).(type) {
	case uint:
		sql = sql.Where(qu.ID.Eq(v))
	case string:
		sql = sql.Where(qu.UserName.Eq(v))
	}

	res, err := sql.First()
	if err != nil {
		var e *Error
		if strings.Contains(err.Error(), "record not found") {
			e = ErrUserNotFound().WithCause(err)
		} else {
			e = ErrServerInternalError().WithCause(err)
		}
		mlog.Log(e.String())
		return nil, e
	}

	return res, nil
}

func CreateUser(user *model.User) *Error {
	err := Q.User.WithContext(context.TODO()).Create(user)
	if err != nil {
		var e *Error
		if strings.Contains(err.Error(), "violates unique constraint") {
			e = ErrUserExist().WithCause(err)
		} else {
			e = ErrServerInternalError().WithCause(err)
		}

		mlog.Log(e.String())
		return e
	}

	return nil
}

func UpdateUser(user *model.User) *Error {
	qu := Q.User
	err := qu.WithContext(context.TODO()).Where(qu.ID.Eq(user.ID)).Save(user)
	if err != nil {
		e := ErrServerInternalError().WithCause(err)
		mlog.Log(e.String())
		return e
	}

	return nil
}

func ListUsers(page api.Pagination) (int64, []*model.User, *Error) {
	qu := Q.User
	sql := qu.WithContext(context.TODO())

	amount, err := sql.Count()
	if err != nil {
		e := ErrServerInternalError().WithCause(err)
		mlog.Log(e.String())
		return 0, nil, e
	}

	res, err := sql.Order(qu.LastLogin.Desc()).Limit(page.Size).Offset((page.Num - 1) * page.Size).Find()
	if err != nil {
		e := ErrServerInternalError().WithCause(err)
		mlog.Log(e.String())
		return 0, nil, e
	}

	return amount, res, nil
}
