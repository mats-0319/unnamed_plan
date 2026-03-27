package api

const URI_CreateNote = "/note/create"

type CreateNoteReq struct {
	IsAnonymous bool   `json:"is_anonymous"`
	Title       string `json:"title"`
	Content     string `json:"content"`
}

type CreateNoteRes struct {
}

const URI_ListNote = "/note/list"

type Note struct {
	NoteID      string `json:"note_id"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	Writer      string `json:"writer"` // writer nickname
	IsAnonymous bool   `json:"is_anonymous"`
	Title       string `json:"title"`
	Content     string `json:"content"`
}

type ListNoteReq struct {
	// 是否仅查询操作员编写的小纸条，需要该值为true且能从请求头解析出操作员用户名才会执行，否则查询全部
	OnlyOperator bool `json:"only_operator"`
	
	Page Pagination `json:"page"`
}

type ListNoteRes struct {
	Amount int64   `json:"amount"`
	Notes  []*Note `json:"notes"`
}

const URI_ModifyNote = "/note/modify"

// ModifyNoteReq can only modify myself note(s)
type ModifyNoteReq struct {
	NoteID      string `json:"note_id"`
	IsAnonymous bool   `json:"is_anonymous"`
	Title       string `json:"title"`
	Content     string `json:"content"`
}

type ModifyNoteRes struct {
}

const URI_DeleteNote = "/note/delete"

// DeleteNoteReq can only delete myself note(s)
type DeleteNoteReq struct {
	NoteID string `json:"note_id"`
}

type DeleteNoteRes struct {
}
