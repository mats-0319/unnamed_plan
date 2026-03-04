package backup

import (
	"github.com/google/uuid"
)

type IBackupRecover interface {
	Backup() error
	Recover() error
}

type BRManager struct {
	list []IBackupRecover
}

func NewBRManager(v ...IBackupRecover) *BRManager {
	return &BRManager{list: v}
}

func (m *BRManager) Backup() error {
	for _, v := range m.list {
		err := v.Backup()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *BRManager) Recover() error {
	for _, v := range m.list {
		err := v.Recover()
		if err != nil {
			return err
		}
	}

	return nil
}

type doBackupRecover[T any] interface {
	Model() T          // *model.User{}
	EmptySlice() []T   // []*model.User{}
	ID(T) uuid.UUID    //
	Update(T)          // update T.backupAt
	Condition() string // "backup_at < updated_at"
	Dir() string       // "user/"

	//DoSomeChangeForTest(T) // 测试用，为了能看出来数据库记录是预设的还是恢复的
}
