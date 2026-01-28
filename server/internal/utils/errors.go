package utils

type ErrorType string

const (
	// change http status code
	ET_ServerInternalError ErrorType = "Server Internal Error"
	ET_UnauthorizedError   ErrorType = "Unauthorized Error"

	ET_ParamsError   ErrorType = "Params Error"
	ET_OperatorError ErrorType = "Operator Error"
)

type ErrorDetail string

const (
	// empty
	ED_Empty ErrorDetail = ""

	// unauthorized
	ED_InvalidAccessToken ErrorDetail = "Invalid Access Token"
	ED_TokenTamperedWith  ErrorDetail = "Access Token has been Tampered With"
	ED_TokenExpired       ErrorDetail = "Token Expired"

	// params
	ED_JsonMarshal   ErrorDetail = "Json Marshal Failed"
	ED_JsonUnmarshal ErrorDetail = "Json Unmarshal Failed"
	ED_HexDecode     ErrorDetail = "Hex Decode Failed"

	// business-user
	ED_NeedAdmin       ErrorDetail = "Need Admin"
	ED_NoChanges       ErrorDetail = "No Changes"
	ED_PwdStructure    ErrorDetail = "Wrong Password Structure"
	ED_PwdVersion      ErrorDetail = "Wrong Password Version"
	ED_PwdParams       ErrorDetail = "Wrong Password Params"
	ED_WrongPwd        ErrorDetail = "Wrong UserName or Password"
	ED_InvalidTotpCode ErrorDetail = "Invalid TOTP Code"
	ED_WrongTotpCode   ErrorDetail = "Wrong TOTP Code"
	ED_InvalidTotpKey  ErrorDetail = "Invalid TOTP Key"
	ED_SamePwd         ErrorDetail = "New Password Can't be Identical to the Old One"

	// business-note
	ED_NeedOwner ErrorDetail = "Not Owner of Data"

	// db
	ED_UserExist    ErrorDetail = "User Already Exist"
	ED_NoteExist    ErrorDetail = "Note Already Exist"
	ED_UserNotFound ErrorDetail = "User Not Found"
	ED_NoteNotFound ErrorDetail = "Note Not Found"
)
