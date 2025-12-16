package db

import (
	"context"
	"strings"

	. "github.com/mats0319/unnamed_plan/server/internal/const"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	api "github.com/mats0319/unnamed_plan/server/internal/http/api/go"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
)

// GetUser query user by 'id'/'username', according to value type
func GetUser[T uint | string](value T) (*model.User, error) {
	qu := dal.Q.User
	sql := qu.WithContext(context.TODO())

	switch v := any(value).(type) {
	case uint:
		sql = sql.Where(qu.ID.Eq(v))
	case string:
		sql = sql.Where(qu.UserName.Eq(v))
	}

	res, err := sql.First()
	if err != nil {
		e := NewError(ET_DBError, ED_Operate).WithCause(err)
		mlog.Log(e.String())
		return nil, e
	}

	return res, nil
}

func CreateUser(user *model.User) error {
	err := dal.Q.User.WithContext(context.TODO()).Create(user)
	if err != nil {
		var e *Error
		if strings.Contains(err.Error(), "violates unique constraint") {
			e = NewError(ET_DBError, ED_DuplicateCreate).WithCause(err)
		} else {
			e = NewError(ET_DBError, ED_Operate).WithCause(err)
		}

		mlog.Log(e.String())
		return e
	}

	return nil
}

func UpdateUser(user *model.User) error {
	qu := dal.Q.User
	err := qu.WithContext(context.TODO()).Where(qu.ID.Eq(user.ID)).Save(user)
	if err != nil {
		e := NewError(ET_DBError, ED_Operate).WithCause(err)
		mlog.Log(e.String())
		return e
	}

	return nil
}

func ListUsers(page api.Pagination) (int64, []*model.User, error) {
	qu := dal.Q.User
	sql := qu.WithContext(context.TODO())

	amount, err := sql.Count()
	if err != nil {
		e := NewError(ET_DBError, ED_Operate).WithCause(err)
		mlog.Log(e.String())
		return 0, nil, err
	}

	res, err := sql.Order(qu.LastLogin.Desc()).Limit(page.Size).Offset((page.Num - 1) * page.Size).Find()
	if err != nil {
		e := NewError(ET_DBError, ED_Operate).WithCause(err)
		mlog.Log(e.String())
		return 0, nil, err
	}

	return amount, res, nil
}
