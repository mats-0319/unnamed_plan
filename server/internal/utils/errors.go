package utils

type ErrorType string

const (
	// change http status code
	ET_BadRequest          ErrorType = "Bad Request"
	ET_Unauthorized        ErrorType = "Unauthorized"
	ET_ServerInternalError ErrorType = "Server Internal Error"

	ET_ParamsError   ErrorType = "Params Error"
	ET_OperatorError ErrorType = "Operator Error"
)

func newError(typ ErrorType, detail ...string) func() *Error {
	return func() *Error { return NewError(typ, detail...) }
}

// 每个变量代表一个函数，一个预设了不同参数的`NewError`函数，使用时需要添加`()`
//
// 使用invalid表示结构/格式错误，使用wrong表示数据错误
var (
	// bad request
	ErrUnregisteredUri = newError(ET_BadRequest, "Unregistered URI")

	// unauthorized
	ErrInvalidAccessToken = newError(ET_Unauthorized, "Invalid Access Token")
	ErrTokenTamperedWith  = newError(ET_Unauthorized, "Access Token has been Tampered With") // hash unmatch
	ErrTokenExpired       = newError(ET_Unauthorized, "Access Token Expired")

	// server internal error
	// 像一个口袋Error，所有不好分类的错误都可以往里装，例如参数全部由服务端代码生成的json marshal error、DB error
	ErrServerInternalError = newError(ET_ServerInternalError)
	ErrDBError             = newError(ET_ServerInternalError, "DB Error")

	// params
	ErrJsonMarshal   = newError(ET_ParamsError, "Json Marshal Failed")
	ErrJsonUnmarshal = newError(ET_ParamsError, "Json Unmarshal Failed")
	ErrHexDecode     = newError(ET_ParamsError, "Hex Decode Failed")

	// business-user
	ErrNeedAdmin       = newError(ET_OperatorError, "Need Admin")
	ErrNoChanges       = newError(ET_ParamsError, "No Changes")
	ErrInvalidPwd      = newError(ET_ParamsError, "Invalid Password Structure")
	ErrPwdParams       = newError(ET_ParamsError, "Wrong Password Params")
	ErrWrongPwd        = newError(ET_ParamsError, "Wrong UserName or Password")
	ErrSamePwd         = newError(ET_ParamsError, "New Password Can't be Identical to the Old One")
	ErrInvalidTotpCode = newError(ET_ParamsError, "Invalid TOTP Code")
	ErrWrongTotpCode   = newError(ET_ParamsError, "Wrong TOTP Code")
	ErrInvalidTotpKey  = newError(ET_ParamsError, "Invalid TOTP Key")

	// business-note
	ErrNeedOwner = newError(ET_OperatorError, "Not Owner of Data")

	// db
	ErrUserExist    = newError(ET_ParamsError, "User Already Exist")
	ErrNoteExist    = newError(ET_ParamsError, "Note Already Exist")
	ErrUserNotFound = newError(ET_ParamsError, "User Not Found")
	ErrNoteNotFound = newError(ET_ParamsError, "Note Not Found")
)
