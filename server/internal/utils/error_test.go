package utils

import (
	"errors"
	"testing"
)

func TestLogStyle(t *testing.T) {
	e := NewError(ET_ParamsError, ED_JsonMarshal).WithCause(errors.New("a new error")).
		WithParam("first param", "first value").
		WithParam("second param", new(string))

	t.Log(e.String())
}
