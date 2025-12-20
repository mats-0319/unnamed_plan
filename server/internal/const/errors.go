package mconst

type ErrorType string

const (
	ET_ServerInternalError ErrorType = "Server Internal Error"
	ET_DBError             ErrorType = "DB Error"
	ET_ParamsError         ErrorType = "Params Error"
	ET_AuthenticateError   ErrorType = "Authentication Error"
)

type ErrorDetail string

const (
	// normal
	ED_UnsupportedURI         ErrorDetail = "Unsupported Request URI"
	ED_UnknownURIOrServerName ErrorDetail = "Unknown URI or Server Name"
	ED_InvalidHttpRequest     ErrorDetail = "Invalid Http Request"
	ED_HttpInvoke             ErrorDetail = "Http Invoke Failed"
	ED_IORead                 ErrorDetail = "IO Read Failed"
	ED_JsonMarshal            ErrorDetail = "Json Marshal Failed"
	ED_JsonUnmarshal          ErrorDetail = "Json Unmarshal Failed"

	// business
	ED_NeedAdmin       ErrorDetail = "Need Admin"
	ED_InvalidPwd      ErrorDetail = "Invalid User Name or Password"
	ED_InvalidTotpCode ErrorDetail = "Invalid Totp Code"
	ED_Base32Decode    ErrorDetail = "Base32 Decode Failed"
	ED_InvalidTotpKey  ErrorDetail = "Invalid Totp Key"
	ED_SamePwd         ErrorDetail = "Can't User Same Password When Re-set"

	// middleware
	ED_InvalidUserIDOrToken ErrorDetail = "Invalid User ID or Token"
	ED_TokenExpired         ErrorDetail = "Token Expired"

	// db
	ED_Operate         ErrorDetail = "DB Operate Failed"
	ED_DuplicateCreate ErrorDetail = "User Already Exist"
)
