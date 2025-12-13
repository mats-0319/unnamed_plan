package db

import (
	"context"

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
		sql = sql.Where(qu.Name.Eq(v))
	}

	res, err := sql.First()
	if err != nil {
		mlog.Log("query user failed", mlog.Field("error", err))
		return nil, err
	}

	return res, nil
}

func CreateUser(user *model.User) error {
	qu := dal.Q.User
	return qu.WithContext(context.TODO()).Create(user)
}

func UpdateUser(user *model.User) error {
	qu := dal.Q.User
	return qu.WithContext(context.TODO()).Where(qu.ID.Eq(user.ID)).Save(user)
}

func ListUsers(page api.Pagination) (int64, []*model.User, error) {
	qu := dal.Q.User
	sql := qu.WithContext(context.TODO())

	amount, err := sql.Count()
	if err != nil {
		mlog.Log("count amount failed", mlog.Field("error", err))
		return 0, nil, err
	}

	res, err := sql.Order(qu.LastLogin.Desc()).Limit(page.Size).Offset((page.Num - 1) * page.Size).Find()
	if err != nil {
		mlog.Log("query user failed", mlog.Field("error", err))
		return 0, nil, err
	}

	return amount, res, nil
}
