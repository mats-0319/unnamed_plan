package test

import (
	"testing"

	"github.com/mats0319/unnamed_plan/server/test/api"
)

func TestApi(t *testing.T) {
	// 测试前需要重置数据库为初始状态（执行建表工具）

	api.UserRegister()
	api.UserLogin()
	api.UserList()
	api.UserModify()
	api.UserAuthenticate()

	api.NoteCreate()
	api.NoteList()
	api.NoteModify()
	api.NoteDelete()
}
