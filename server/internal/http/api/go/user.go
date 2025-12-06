package api

const URI_Login = "/login"

type LoginReq struct {
	UserName string `json:"user_name"`
	Password string `json:"password"` // sha256('text')
	TotpCode string `json:"totp_code"`
}

type LoginRes struct {
	Res   ResBase `json:"res"`
	Token string  `json:"token"`
}

const URI_ListUser = "/user/list"

type User struct {
	ID         uint   `json:"id"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
	Name       string `json:"name"`     // login name
	Nickname   string `json:"nickname"` // display name
	TotpKey    string `json:":16"`      // 允许为空，需要设置后启动
	IsLocked   bool   `json:"is_locked"`
	Permission uint8  `json:"permission"` // 权限等级
	Creator    string `json:"creator"`    // 创建人，user.name
	LastLogin  int64  `json:"last_login"` // timestamp, unit: milli
}

// ListUserReq 考虑添加按照字段查询、按照字段排序
type ListUserReq struct {
	Operator  uint64     `json:"operator"`
	Page      Pagination `json:"page"`
	SeeLocked bool       `json:"see_locked"`
}

type ListUserRes struct {
	Res        ResBase `json:"res"`
	UserAmount int     `json:"user_amount"` // 用户表总数
	ListAmount int     `json:"list_amount"` // 符合查询条件的用户数
	Users      []User  `json:"users"`
}

const URI_CreateUser = "/user/create"

type CreateUserReq struct {
	Operator   uint64 `json:"operator"`  // find user by id, user.name is 'creator'
	UserName   string `json:"user_name"` // nickname is same, user can modify later
	Password   string `json:"password"`  // sha256('text'), server generate 'salt' and save it
	Permission uint8  `json:"permission"`
}

type CreateUserRes struct {
	Res ResBase `json:"res"`
}

const URI_LockUser = "/user/lock"

type LockUserReq struct {
	Operator uint64 `json:"operator"`
}

type LockUserRes struct {
	Res ResBase `json:"res"`
}

const URI_UnlockUser = "/user/unlock"

type UnlockUserReq struct {
	Operator uint64 `json:"operator"`
}

type UnlockUserRes struct {
	Res ResBase `json:"res"`
}

const URI_ModifyInfo = "/user/modify-info"

// ModifyInfoReq 属性字段为空，视作不修改对应字段，有专属的bool变量标识是否修改的字段不适用该默认规则
type ModifyInfoReq struct {
	Operator     uint64 `json:"operator"`
	Nickname     string `json:"nickname"`
	Password     string `json:"password"`       // sha256('text'), check: can't be same between old and new pwd
	ModifyTkFlag bool   `json:"modify_tk_flag"` // if modify totp key
	TotpKey      string `json:"totp_key"`       // length: 16
}

type ModifyInfoRes struct {
	Res ResBase `json:"res"`
}

const URI_AdjustPermission = "/user/adjust-permission"

type AdjustPermissionReq struct {
	Operator   uint64 `json:"operator"`
	UserID     uint64 `json:"user_id"`
	Permission uint8  `json:"permission"`
}

type AdjustPermissionRes struct {
	Res ResBase `json:"res"`
}

const URI_Authenticate = "/user/authenticate"

type AuthenticateReq struct {
	UserID   uint64 `json:"user_id"`
	Password string `json:"password"` // sha256('text')
}

type AuthenticateRes struct {
	Res ResBase `json:"res"`
}
