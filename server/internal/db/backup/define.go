package backup

import (
	"github.com/google/uuid"
)

type IBackupRecover interface {
	Backup()
	Recover()
}

type BRManager struct {
	List []IBackupRecover
}

func NewBRManager(v ...IBackupRecover) *BRManager {
	return &BRManager{List: v}
}

func (m *BRManager) Backup() {
	for _, v := range m.List {
		v.Backup()
	}
}

func (m *BRManager) Recover() {
	for _, v := range m.List {
		v.Recover()
	}
}

// 我们提供的备份/恢复方法基于该接口类型的变量；反过来说，如果你实现了该接口，就可以使用我们提供的备份/恢复方法
// 建议类型T为指针类型的数据库结构体
// 我们在程序中硬编码了一些参数（例如分页大小、单表最大文件数量、默认备份软删除的数据等），我们认为这是可以接受的
type doBackupRecover[T any] interface {
	Model() T              // *model.User{}
	EmptySlice() []T       // []*model.User{}
	ID(T) uuid.UUID        //
	Update(T, int64)       // 备份时使用，维护 T.backupAt
	ColumnNames() []string // 恢复时使用，禁用updateAt字段的自动更新（备份时使用gorm提供的函数）
	Condition() string     // 备份条件 "backup_at < updated_at"
	Dir() string           // "user/"

	//DoSomeChangesForTest([]T) // 测试用，为了能看出来数据库记录是不是恢复的
}
