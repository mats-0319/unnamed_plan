package utils

type ErrorType string

const (
	// change http status code
	ET_ServerInternalError ErrorType = "Server Internal Error"
	ET_UnauthorizedError   ErrorType = "Unauthorized Error"

	ET_ParamsError   ErrorType = "Params Error"
	ET_OperatorError ErrorType = "Operator Error"
)

func newError(typ ErrorType, detail ...string) func() *Error {
	return func() *Error { return NewError(typ, detail...) }
}

// 每个变量代表一个函数，一个预设了不同参数的`NewError`函数，使用时需要添加`()`
// 其实严格来说，应该把internal用到的error和cmd用到的error也区分开，但我懒得搞了
var (
	// server internal error
	ErrServerInternalError = newError(ET_ServerInternalError)

	// unauthorized error
	ErrInvalidAccessToken = newError(ET_UnauthorizedError, "Invalid Access Token")
	ErrTokenTamperedWith  = newError(ET_UnauthorizedError, "Access Token has been Tampered With")
	ErrTokenExpired       = newError(ET_UnauthorizedError, "Token Expired")

	// params
	ErrJsonMarshal   = newError(ET_ParamsError, "Json Marshal Failed")
	ErrJsonUnmarshal = newError(ET_ParamsError, "Json Unmarshal Failed")
	ErrHexDecode     = newError(ET_ParamsError, "Hex Decode Failed")

	// business-user
	ErrNeedAdmin       = newError(ET_OperatorError, "Need Admin")
	ErrNoChanges       = newError(ET_ParamsError, "No Changes")
	ErrPwdStructure    = newError(ET_ParamsError, "Wrong Password Structure")
	ErrPwdVersion      = newError(ET_ParamsError, "Wrong Password Version")
	ErrPwdParams       = newError(ET_ParamsError, "Wrong Password Params")
	ErrWrongPwd        = newError(ET_ParamsError, "Wrong UserName or Password")
	ErrInvalidTotpCode = newError(ET_ParamsError, "Invalid TOTP Code")
	ErrWrongTotpCode   = newError(ET_ParamsError, "Wrong TOTP Code")
	ErrInvalidTotpKey  = newError(ET_ParamsError, "Invalid TOTP Key")
	ErrSamePwd         = newError(ET_ParamsError, "New Password Can't be Identical to the Old One")

	// business-note
	ErrNeedOwner = newError(ET_OperatorError, "Not Owner of Data")

	// db
	ErrUserExist    = newError(ET_ParamsError, "User Already Exist")
	ErrNoteExist    = newError(ET_ParamsError, "Note Already Exist")
	ErrUserNotFound = newError(ET_ParamsError, "User Not Found")
	ErrNoteNotFound = newError(ET_ParamsError, "Note Not Found")
)
