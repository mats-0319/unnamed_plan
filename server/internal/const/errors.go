package mconst

type ErrorType string

const (
	// change http status code
	ET_ServerInternalError ErrorType = "Server Internal Error"
	ET_UnauthorizedError   ErrorType = "Authentication Error"

	ET_ParamsError   ErrorType = "Params Error"
	ET_OperatorError ErrorType = "Operator Error"
)

type ErrorDetail string

const (
	// empty
	ED_Empty ErrorDetail = ""

	// unauthorized
	ED_InvalidUserIDOrToken ErrorDetail = "Invalid User ID or Token"
	ED_TokenExpired         ErrorDetail = "Token Expired"

	// params
	ED_JsonMarshal   ErrorDetail = "Json Marshal Failed"
	ED_JsonUnmarshal ErrorDetail = "Json Unmarshal Failed"

	// business-user
	ED_NeedAdmin       ErrorDetail = "Need Admin"
	ED_InvalidPwd      ErrorDetail = "Invalid UserName or Password"
	ED_InvalidTotpCode ErrorDetail = "Invalid TOTP Code"
	ED_InvalidTotpKey  ErrorDetail = "Invalid TOTP Key"
	ED_SamePwd         ErrorDetail = "Can't User Same Password When Re-set"

	// business-note
	ED_NeedOwner     ErrorDetail = "Not Owner of Data"
	ED_ModifyNothing ErrorDetail = "No Changes"

	// db
	ED_UserExist ErrorDetail = "User Already Exist"
	ED_NoteExist ErrorDetail = "Note Already Exist"
)
