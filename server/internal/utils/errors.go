package utils

// 每个变量代表一个函数，一个预设了不同参数的`NewError`函数，使用时需要添加`()`
// 错误码为5位十进制数，第一位表示错误类型，2～3位表示功能模块，4～5位表示错误序号
// 使用invalid表示结构错误，使用wrong表示数值错误
var (
	// change http code error (code: start with 'http code')
	ErrInvalidUri             = newError(400, 40001, "Invalid URI")
	ErrInvalidAccessToken     = newError(401, 40101, "Invalid Access Token")
	ErrWrongAccessTokenHash   = newError(401, 40102, "Wrong Access Token Hash")
	ErrDecodeAccessToken      = newError(401, 40103, "Decode Access Token Failed")
	ErrDeserializeAccessToken = newError(401, 40104, "Deserialize Access Token Failed")
	ErrAccessTokenExpired     = newError(401, 40105, "Access Token Expired")
	ErrInvalidMfaToken        = newError(401, 40106, "Invalid MFA Token")
	ErrDecodeMfaToken         = newError(401, 40107, "Decode MFA Token Failed")
	ErrDeserializeMfaToken    = newError(401, 40108, "Deserialize MFA Token Failed")
	ErrNoMfaToken             = newError(401, 40109, "Input User Name and Password First")
	ErrWrongMfaToken          = newError(401, 40110, "Wrong MFA Token")
	ErrMfaTokenExpired        = newError(401, 40111, "MFA Token Expired")
	ErrServerInternalError    = newError(500, 50001, "Sever Internal Error")
	ErrDBError                = newError(500, 50002, "DB Error")

	// params error (1)
	// general (00) / db (01)
	ErrDeserializeHttpReqParam = newBusinessError(10001, "Deserialize Http Request Param Failed")
	ErrInvalidParams           = newBusinessError(10002, "Invalid Params")
	ErrUserExist               = newBusinessError(10101, "User Already Exist")
	ErrNoteExist               = newBusinessError(10102, "Note Already Exist")
	ErrUserNotFound            = newBusinessError(10103, "User Not Found")
	ErrNoteNotFound            = newBusinessError(10104, "Note Not Found")

	// business error (2)
	// general (00) / user (01) / note (02) / game score (03)
	ErrNoChanges        = newBusinessError(20001, "No Changes")
	ErrPermissionDenied = newBusinessError(20002, "Permission Denied")
	ErrInvalidPassword  = newBusinessError(20101, "Invalid Password")
	ErrInvalidPwdSalt   = newBusinessError(20102, "Invalid Password Salt")
	ErrInvalidPwdKey    = newBusinessError(20103, "Invalid Password Key")
	ErrWrongPassword    = newBusinessError(20104, "Wrong Password")
	ErrInvalidTotpCode  = newBusinessError(20105, "Invalid TOTP Code")
	ErrInvalidTotpKey   = newBusinessError(20106, "Invalid TOTP Key")
	ErrWrongTotpCode    = newBusinessError(20107, "Wrong TOTP Code")
	ErrSamePassword     = newBusinessError(20108, "New Password Can't be Identical to the Old One")
	ErrInvalidGameName  = newBusinessError(20301, "Invalid Game Name")
)

func newError(httpCode int, code int, detail string) func() *Error {
	return func() *Error { return NewError(httpCode, code, detail) }
}

func newBusinessError(code int, detail string) func() *Error {
	return newError(200, code, detail)
}
