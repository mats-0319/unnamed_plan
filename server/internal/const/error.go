package mconst

import (
	"fmt"
	"runtime"
	"strings"
)

type Error struct {
	Typ    ErrorType
	Detail ErrorDetail
	Cause  error
	Params map[string]any
	Stack  []uintptr
}

var _ error = (*Error)(nil)

func NewError(typ ErrorType, detail ...ErrorDetail) *Error {
	errDetail := ED_Empty
	if len(detail) > 0 {
		errDetail = detail[0]
	}

	var stack [32]uintptr
	n := runtime.Callers(2, stack[:]) // skip 'runtime.caller' and 'NewError'

	return &Error{
		Typ:    typ,
		Detail: errDetail,
		Params: make(map[string]any),
		Stack:  stack[:n],
	}
}

// Error simple return to web
func (e *Error) Error() string {
	detailStr := ""
	if len(e.Detail) > 0 {
		detailStr = fmt.Sprintf(", detail: %s", e.Detail)
	}

	return fmt.Sprintf("error type: %s%s", e.Typ, detailStr)
}

// String print all details, use in dev
func (e *Error) String() string {
	detailStr := ""
	if len(e.Detail) > 0 {
		detailStr = fmt.Sprintf(", detail: %s", e.Detail)
	}

	errStr := ""
	if e.Cause != nil {
		errStr = fmt.Sprintf("\nerr: %v", e.Cause)
	}

	paramStr := ""
	if len(e.Params) > 0 {
		paramStr = fmt.Sprintf("\nparams: %#v", e.Params)
	}

	return fmt.Sprintf("error type: %s%s%s%s\nstack trace: \n%s\n",
		e.Typ, detailStr, errStr, paramStr, e.stackTrace())
}

func (e *Error) WithCause(err error) *Error {
	e.Cause = err
	return e
}

func (e *Error) WithParam(key string, value any) *Error {
	if e.Params == nil {
		e.Params = make(map[string]any)
	}

	e.Params[key] = value
	return e
}

func (e *Error) stackTrace() string {
	var builder strings.Builder
	frames := runtime.CallersFrames(e.Stack)

	for {
		frame, more := frames.Next()
		builder.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}

	return builder.String()
}
