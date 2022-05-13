package e_test

import (
	"bytes"
	"e"
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	msg := "primary"
	ctx := map[string]interface{}{
		"foo": "bar",
	}

	err := e.New(msg, ctx)

	if !e.Is(err) {
		t.Error("failed type check")
	}

	if err.Error() != msg {
		t.Errorf("failed message check '%s' != '%s'", err.Error(), msg)
	}

	if !reflect.DeepEqual(e.Context(err), ctx) {
		t.Error("failed context check")
	}

	if e.Depth(err) != 1 {
		t.Errorf("failed depth check '%d' != '%d'", e.Depth(err), 1)
	}

	if e.Primary(err) != err {
		t.Errorf("failed primary check")
	}

	if e.Cause(err) != nil {
		t.Errorf("failed cause check")
	}

	if e.Root(err) != nil {
		t.Errorf("failed root check")
	}

	buf, _ := e.JSON(err)

	if !bytes.Equal(buf, []byte(`[{"message":"primary","context":{"foo":"bar"}}]`)) {
		t.Errorf("failed json check")
	}
}

func TestWrap(t *testing.T) {
	rt := fmt.Errorf("root")
	msg := "primary"
	ctx := map[string]interface{}{
		"foo": "bar",
	}

	pri := e.Wrap(rt, msg, ctx)

	if !e.Is(pri) {
		t.Error("failed type check")
	}

	if pri.Error() != msg {
		t.Errorf("failed message check '%s' != '%s'", pri.Error(), msg)
	}

	if !reflect.DeepEqual(e.Context(pri), ctx) {
		t.Error("failed context check")
	}

	if e.Depth(pri) != 2 {
		t.Errorf("failed depth check '%d' != '%d'", e.Depth(pri), 2)
	}

	if e.Primary(pri) != pri {
		t.Errorf("failed primary check")
	}

	if e.Cause(pri) != rt {
		t.Errorf("failed cause check")
	}

	if e.Root(pri) != rt {
		t.Errorf("failed root check")
	}

	buf, _ := e.JSON(pri)

	if !bytes.Equal(buf, []byte(`[{"message":"primary","context":{"foo":"bar"}},{"message":"root","context":null}]`)) {
		t.Errorf("failed json check")
	}

	msg = "secondary"
	ctx = map[string]interface{}{
		"hot": "dog",
	}

	sec := e.Wrap(pri, msg, ctx)

	if !e.Is(sec) {
		t.Error("failed type check")
	}

	if sec.Error() != msg {
		t.Errorf("failed message check '%s' != '%s'", sec.Error(), msg)
	}

	if !reflect.DeepEqual(e.Context(sec), ctx) {
		t.Error("failed context check")
	}

	if e.Depth(sec) != 3 {
		t.Errorf("failed depth check '%d' != '%d'", e.Depth(sec), 3)
	}

	if e.Primary(sec) != pri {
		t.Errorf("failed primary check")
	}

	if e.Cause(sec) != pri {
		t.Errorf("failed cause check")
	}

	if e.Root(sec) != rt {
		t.Errorf("failed root check")
	}

	buf, _ = e.JSON(sec)

	if !bytes.Equal(buf, []byte(`[{"message":"secondary","context":{"hot":"dog"}},{"message":"primary","context":{"foo":"bar"}},{"message":"root","context":null}]`)) {
		t.Errorf("failed json check")
	}
}
