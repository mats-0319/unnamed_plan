package api

const URI_Register = "/register"

type RegisterReq struct {
	UserName string `json:"user_name"` // unique, also default nickname
	Password string `json:"password"`  // hex(sha256('text'))
}

type RegisterRes struct {
}

const URI_Login = "/login"

type LoginReq struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type LoginRes struct {
	UserName  string `json:"user_name"`
	Nickname  string `json:"nickname"`
	IsAdmin   bool   `json:"is_admin"`
	EnableMFA bool   `json:"enable_mfa"`
	MFAToken  string `json:"mfa_token"` // empty when disable 2fa
}

const URI_LoginMFA = "/login-mfa"

type LoginMFAReq struct {
	MFAToken string `json:"mfa_token"`
	TOTPCode string `json:"totp_code"`
}

type LoginMFARes struct {
	UserName  string `json:"user_name"`
	Nickname  string `json:"nickname"`
	IsAdmin   bool   `json:"is_admin"`
	EnableMFA bool   `json:"enable_mfa"`
}

const URI_ListUser = "/user/list"

type User struct {
	UserName  string `json:"user_name"` // login name, can't modify
	Nickname  string `json:"nickname"`  // display name
	CreatedAt int64  `json:"created_at"`
	IsAdmin   bool   `json:"is_admin"`
	EnableMFA bool   `json:"enable_mfa"`
	LastLogin int64  `json:"last_login"` // timestamp, unit: milli
}

type ListUserReq struct {
	Page Pagination `json:"page"`
}

type ListUserRes struct {
	Count int64   `json:"count"` // 符合查询条件的用户总数
	Users []*User `json:"users"`
}

const URI_ModifyUser = "/user/modify"

// ModifyUserReq string类型的属性为空，视为不修改对应字段
type ModifyUserReq struct {
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	EnableMFA bool   `json:"enable_mfa"`
	TOTPKey   string `json:"totp_key"`
}

type ModifyUserRes struct {
}
