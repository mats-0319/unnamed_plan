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

func NewError(typ ErrorType, detail ErrorDetail) *Error {
	var stack [32]uintptr
	n := runtime.Callers(2, stack[:]) // skip 'caller' and 'new'

	return &Error{
		Typ:    typ,
		Detail: detail,
		Params: make(map[string]any),
		Stack:  stack[:n],
	}
}

// Error simple return to web
func (e *Error) Error() string {
	return fmt.Sprintf("error type: %s, detail: %s", e.Typ, e.Detail)
}

// String print all details, use in dev
func (e *Error) String() string {
	errStr := ""
	if e.Cause != nil {
		errStr = fmt.Sprintf(", err: %v\n", e.Cause)
	}

	paramStr := ""
	if len(e.Params) > 0 {
		paramStr = fmt.Sprintf("params: %+v\n", e.Params)
	}

	return fmt.Sprintf("error type: %s, detail: %s%s%sstack trace: %s\n",
		e.Typ, e.Detail, errStr, paramStr, e.stackTrace())
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
