package backup

import (
	"time"

	"github.com/google/uuid"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
)

type UserBR struct {
}

var _ IBackupRecover = (*UserBR)(nil)

func (i *UserBR) Backup() error {
	return Backup[*model.User](i) // 这里会隐式验证类型是否实现了`doBackupRecover`接口，不需要显式验证
}

func (i *UserBR) Recover() error {
	return Recover[*model.User](i)
}

func (i *UserBR) Model() *model.User {
	return &model.User{}
}

func (i *UserBR) EmptySlice() []*model.User {
	return make([]*model.User, 0)
}

func (i *UserBR) ID(t *model.User) uuid.UUID {
	return t.ID
}

func (i *UserBR) Update(data *model.User) {
	data.BackupAt = time.Now().Add(5 * time.Minute).UnixMilli()
}

func (i *UserBR) Condition() string {
	return "backup_at < updated_at"
}

func (i *UserBR) Dir() string {
	return "user/"
}

func (i *UserBR) DoSomeChangeForTest(data *model.User) {
	data.Nickname = "Recover - " + data.Nickname
}

type NoteBR struct {
}

var _ IBackupRecover = (*NoteBR)(nil)

func (i *NoteBR) Backup() error {
	return Backup[*model.Note](i)
}

func (i *NoteBR) Recover() error {
	return Recover[*model.Note](i)
}

func (i *NoteBR) Model() *model.Note {
	return &model.Note{}
}

func (i *NoteBR) EmptySlice() []*model.Note {
	return make([]*model.Note, 0)
}

func (i *NoteBR) ID(t *model.Note) uuid.UUID {
	return t.ID
}

func (i *NoteBR) Update(data *model.Note) {
	data.BackupAt = time.Now().Add(5 * time.Minute).UnixMilli()
}

func (i *NoteBR) Condition() string {
	return "backup_at < updated_at"
}

func (i *NoteBR) Dir() string {
	return "note/"
}

func (i *NoteBR) DoSomeChangeForTest(data *model.Note) {
	data.WriterName = "Recover - " + data.WriterName
}
