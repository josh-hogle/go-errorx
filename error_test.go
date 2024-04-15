package errorx_test

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"testing"
	"time"

	"go.joshhogle.dev/errorx"
)

// TODO: implement testing and benchmarks

type GenericError struct {
	*errorx.BaseError
}

func NewGenericError(code int, err error) *GenericError {
	return &GenericError{
		BaseError: errorx.NewBaseError(code, err),
	}
}

func (e *GenericError) Error() string {
	buf := bytes.NewBuffer(nil)
	attrs := e.Attrs()
	fmt.Fprintf(buf, "a generic error has occurred: %s (code=%d)", e.InternalError().Error(), e.Code())
	if len(attrs) > 0 {
		buf.WriteString(" [")
		for k, v := range attrs {
			fmt.Fprintf(buf, " %s=%v", k, v)
		}
		buf.WriteString(" ]")
	}
	for _, n := range e.NestedErrors() {
		buf.WriteString("\n\t")
		buf.WriteString(n.Error())
	}
	return buf.String()
}

func TestBaseError1(t *testing.T) {
	err := errors.New("this is an error")
	e1 := NewGenericError(100, err)
	t.Logf("e1: %s\n", e1.Error())
	e2 := NewGenericError(101, err)
	e2.Append(e1)
	t.Logf("e2: %s\n", e2.Error())
	e3 := NewGenericError(102, err)
	e3.Append(e2)
	e3.WithAttrs(map[string]any{
		"key1": "value1",
		"key2": 2334,
		"key3": time.Now().UTC(),
	})
	t.Logf("e3: %s\n", e3.Error())

	_, f, _, _ := runtime.Caller(0)
	t.Logf("file: %s", f)
}
