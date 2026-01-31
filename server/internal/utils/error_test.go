package utils

import (
	"errors"
	"testing"
)

func TestLogStyle(t *testing.T) {
	e := ErrJsonMarshal().WithCause(errors.New("a new error")).
		WithParam("first param", "first value").
		WithParam("second param", new(string))

	t.Log(e.String())
}

func TestReuse(t *testing.T) {
	e := ErrServerInternalError().WithParam("key", "value")
	t.Log(e.String())
	
	e = ErrServerInternalError()
	t.Log(e.String())
}
