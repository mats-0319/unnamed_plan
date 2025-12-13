package mconst

import (
	"fmt"
)

type Error struct {
	Typ    ErrorType
	Detail ErrorDetail
	Params string
}

func NewError(typ ErrorType, detail ErrorDetail, params string) error {
	return &Error{
		Typ:    typ,
		Detail: detail,
		Params: params,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("error type: %s, detail: %s, param(s): %s", e.Typ, e.Detail, e.Params)
}

type ErrorType string

const (
	ET_InitError           ErrorType = "Init Error"
	ET_ServerInternalError ErrorType = "Server Internal Error" // 程序内部逻辑错误
	ET_DBError             ErrorType = "DB Error"
	ET_HttpError           ErrorType = "Http Error"
)

type ErrorDetail string

const (
	ED_UnsupportedURI ErrorDetail = "Unsupported Request URI"
	ED_ExecMiddleware ErrorDetail = "Exec Middleware Failed"
)
