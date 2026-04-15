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
	Enable2FA bool   `json:"enable_2fa"`
	MfaToken  string `json:"mfa_token"` // empty if disable 2fa
}

const URI_LoginTotp = "/login-totp"

type LoginTotpReq struct {
	UserName string `json:"user_name"`
	TotpCode string `json:"totp_code"`
}

type LoginTotpRes struct {
	UserName  string `json:"user_name"`
	Nickname  string `json:"nickname"`
	IsAdmin   bool   `json:"is_admin"`
	Enable2FA bool   `json:"enable_2fa"`
}

const URI_ListUser = "/user/list"

type User struct {
	UserName  string `json:"user_name"` // login name, can't modify
	Nickname  string `json:"nickname"`  // display name
	CreatedAt int64  `json:"created_at"`
	IsAdmin   bool   `json:"is_admin"`
	Enable2FA bool   `json:"enable_2fa"`
	LastLogin int64  `json:"last_login"` // timestamp, unit: milli
}

type ListUserReq struct {
	Page Pagination `json:"page"`
}

type ListUserRes struct {
	Amount int64   `json:"amount"` // 符合查询条件的用户总数
	Users  []*User `json:"users"`
}

const URI_ModifyUser = "/user/modify"

// ModifyUserReq string类型的属性为空，视为不修改对应字段
type ModifyUserReq struct {
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	Enable2FA bool   `json:"enable_2fa"`
	TotpKey   string `json:"totp_key"`
}

type ModifyUserRes struct {
}
