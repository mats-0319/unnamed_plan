package backup

import (
	"github.com/google/uuid"
	"github.com/mats0319/unnamed_plan/server/internal/db/dal"
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
)

type UserBR struct {
}

var _ IBackupRecover = (*UserBR)(nil)

//var _ doBackupRecover[*model.User] = (*UserBR)(nil)

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

func (i *UserBR) Update(data *model.User, timestamp int64) {
	data.BackupAt = timestamp
}

func (i *UserBR) ColumnNames() []string {
	return []string{
		dal.User.ID.ColumnName().String(),
		dal.User.CreatedAt.ColumnName().String(),
		dal.User.UpdatedAt.ColumnName().String(),
		dal.User.BackupAt.ColumnName().String(),
		dal.User.UserName.ColumnName().String(),
		dal.User.Nickname.ColumnName().String(),
		dal.User.Password.ColumnName().String(),
		dal.User.IsAdmin.ColumnName().String(),
		dal.User.Enable2FA.ColumnName().String(),
		dal.User.TotpKey.ColumnName().String(),
	}
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

//var _ doBackupRecover[*model.Note] = (*NoteBR)(nil)

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

func (i *NoteBR) Update(data *model.Note, timestamp int64) {
	data.BackupAt = timestamp
}

func (i *NoteBR) ColumnNames() []string {
	return []string{
		dal.Note.ID.ColumnName().String(),
		dal.Note.CreatedAt.ColumnName().String(),
		dal.Note.UpdatedAt.ColumnName().String(),
		dal.Note.BackupAt.ColumnName().String(),
		dal.Note.NoteID.ColumnName().String(),
		dal.Note.Writer.ColumnName().String(),
		dal.Note.WriterName.ColumnName().String(),
		dal.Note.IsAnonymous.ColumnName().String(),
		dal.Note.Title.ColumnName().String(),
		dal.Note.Content.ColumnName().String(),
	}
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
