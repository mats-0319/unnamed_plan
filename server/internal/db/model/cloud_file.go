package model

type CloudFile struct {
	Model
	FileID    string `gorm:"size:64;unique;not null"` // sha256('user id'+timestamp)
	Name      string `gorm:"not null"`                // when user upload
	Extension string // format: 'xxx' e.g. 'jpg'/'pdf', no '.'(dot)
	Uploader  string `gorm:"not null"` // user.name
	Size      int64
	Hash      string `gorm:"size:64;not null"` // sha1('file')
	IsDeleted bool
}
