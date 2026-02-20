package utils

import (
	"errors"
	"testing"
)

func TestLogStyle(t *testing.T) {
	e := ErrServerInternalError().WithCause(errors.New("a new error")).
		WithParam("first param", "first value").
		WithParam("second param", 10000)

	t.Log(e.String())
}
