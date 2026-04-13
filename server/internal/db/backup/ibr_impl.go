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

func (u *UserBR) Backup() {
	Backup[*model.User](u) // 这里会隐式验证类型是否实现了`doBackupRecover`接口，不需要显式验证
}

func (u *UserBR) Recover() {
	Recover[*model.User](u)
}

func (u *UserBR) Model() *model.User {
	return &model.User{}
}

func (u *UserBR) EmptySlice() []*model.User {
	return make([]*model.User, 0)
}

func (u *UserBR) ID(t *model.User) uuid.UUID {
	return t.ID
}

func (u *UserBR) Update(data *model.User, timestamp int64) {
	data.BackupAt = timestamp
}

func (u *UserBR) ColumnNames() []string {
	return []string{
		dal.User.ID.ColumnName().String(),
		dal.User.CreatedAt.ColumnName().String(),
		dal.User.UpdatedAt.ColumnName().String(),
		dal.User.DeletedAt.ColumnName().String(),
		dal.User.BackupAt.ColumnName().String(),
		dal.User.UserName.ColumnName().String(),
		dal.User.Nickname.ColumnName().String(),
		dal.User.Password.ColumnName().String(),
		dal.User.IsAdmin.ColumnName().String(),
		dal.User.Enable2FA.ColumnName().String(),
		dal.User.TotpKey.ColumnName().String(),
	}
}

func (u *UserBR) Condition() string {
	return "backup_at < updated_at"
}

func (u *UserBR) Dir() string {
	return "user/"
}

func (u *UserBR) DoSomeChangesForTest(users []*model.User) {
	for _, user := range users {
		user.Nickname = "Recover - " + user.Nickname
	}
}

type NoteBR struct {
}

var _ IBackupRecover = (*NoteBR)(nil)

//var _ doBackupRecover[*model.Note] = (*NoteBR)(nil)

func (n *NoteBR) Backup() {
	Backup[*model.Note](n)
}

func (n *NoteBR) Recover() {
	Recover[*model.Note](n)
}

func (n *NoteBR) Model() *model.Note {
	return &model.Note{}
}

func (n *NoteBR) EmptySlice() []*model.Note {
	return make([]*model.Note, 0)
}

func (n *NoteBR) ID(t *model.Note) uuid.UUID {
	return t.ID
}

func (n *NoteBR) Update(data *model.Note, timestamp int64) {
	data.BackupAt = timestamp
}

func (n *NoteBR) ColumnNames() []string {
	return []string{
		dal.Note.ID.ColumnName().String(),
		dal.Note.CreatedAt.ColumnName().String(),
		dal.Note.UpdatedAt.ColumnName().String(),
		dal.Note.DeletedAt.ColumnName().String(),
		dal.Note.BackupAt.ColumnName().String(),
		dal.Note.NoteID.ColumnName().String(),
		dal.Note.Writer.ColumnName().String(),
		dal.Note.WriterName.ColumnName().String(),
		dal.Note.IsAnonymous.ColumnName().String(),
		dal.Note.Title.ColumnName().String(),
		dal.Note.Content.ColumnName().String(),
	}
}

func (n *NoteBR) Condition() string {
	return "backup_at < updated_at"
}

func (n *NoteBR) Dir() string {
	return "note/"
}

func (n *NoteBR) DoSomeChangesForTest(notes []*model.Note) {
	for _, note := range notes {
		note.WriterName = "Recover - " + note.WriterName
	}
}

type FlipGameScore struct{}

var _ IBackupRecover = (*FlipGameScore)(nil)

//var _ doBackupRecover[*model.FlipGameScore] = (*FlipGameScore)(nil)

func (f *FlipGameScore) Backup() {
	Backup[*model.FlipGameScore](f)
}

func (f *FlipGameScore) Recover() {
	Recover[*model.FlipGameScore](f)
}

func (f *FlipGameScore) Model() *model.FlipGameScore {
	return &model.FlipGameScore{}
}

func (f *FlipGameScore) EmptySlice() []*model.FlipGameScore {
	return make([]*model.FlipGameScore, 0)
}

func (f *FlipGameScore) ID(t *model.FlipGameScore) uuid.UUID {
	return t.ID
}

func (f *FlipGameScore) Update(data *model.FlipGameScore, timestamp int64) {
	data.BackupAt = timestamp
}

func (f *FlipGameScore) ColumnNames() []string {
	return []string{
		dal.FlipGameScore.ID.ColumnName().String(),
		dal.FlipGameScore.CreatedAt.ColumnName().String(),
		dal.FlipGameScore.UpdatedAt.ColumnName().String(),
		dal.FlipGameScore.DeletedAt.ColumnName().String(),
		dal.FlipGameScore.BackupAt.ColumnName().String(),
		dal.FlipGameScore.Score.ColumnName().String(),
		dal.FlipGameScore.Result.ColumnName().String(),
		dal.FlipGameScore.Player.ColumnName().String(),
		dal.FlipGameScore.PlayerName.ColumnName().String(),
	}
}

func (f *FlipGameScore) Condition() string {
	return "backup_at < updated_at"
}

func (f *FlipGameScore) Dir() string {
	return "flip_game_score/"
}

func (f *FlipGameScore) DoSomeChangesForTest(scores []*model.FlipGameScore) {
	for _, score := range scores {
		score.PlayerName = "Recover - " + score.PlayerName
	}
}
