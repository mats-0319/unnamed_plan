package dal

import (
	"context"
	"strings"

	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	mlog "github.com/mats0319/unnamed_plan/server/internal/log"
	. "github.com/mats0319/unnamed_plan/server/internal/utils"
	"gorm.io/gorm"
)

func GetUser(userName string) (*model.User, *Error) {
	res, err := User.WithContext(context.TODO()).Where(User.UserName.Eq(userName)).First()
	if err != nil {
		var e *Error
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			e = ErrUserNotFound().WithCause(err)
		} else {
			e = ErrDBError().WithCause(err)
		}
		mlog.Error(e.String())
		return nil, e
	}

	return res, nil
}

func CreateUser(user *model.User) *Error {
	if err := User.WithContext(context.TODO()).Create(user); err != nil {
		var e *Error
		if strings.Contains(err.Error(), "violates unique constraint") {
			e = ErrUserExist().WithCause(err)
		} else {
			e = ErrDBError().WithCause(err)
		}

		mlog.Error(e.String())
		return e
	}

	return nil
}

func UpdateUser(user *model.User) *Error {
	if err := User.WithContext(context.TODO()).Where(User.ID.Eq(user.ID)).Save(user); err != nil {
		e := ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return e
	}

	return nil
}

func ListUser(pageSize int, pageNum int) (int64, []*model.User, *Error) {
	sql := User.WithContext(context.TODO())

	amount, err := sql.Count()
	if err != nil {
		e := ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return 0, nil, e
	}

	users, err := sql.Order(User.UpdatedAt.Desc()).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find()
	if err != nil {
		e := ErrDBError().WithCause(err)
		mlog.Error(e.String())
		return 0, nil, e
	}

	return amount, users, nil
}
