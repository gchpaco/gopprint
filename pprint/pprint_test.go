package pprint

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Output(elt Element, width int) (string, error) {
	buffer := new(bytes.Buffer)

	err := PrettyPrint(elt, width, buffer)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func TestSimpleText(t *testing.T) {
	handle := NewText("Some text")

	out, err := Output(handle, 10)
	if assert.NoError(t, err) {
		assert.Equal(t, "Some text", out)
	}
	out, err = Output(handle, 4)
	if assert.NoError(t, err) {
		assert.Equal(t, "Some text", out)
	}
}

func TestSimpleCond(t *testing.T) {
	handle := NewConcat(NewText("Some text"), CondLB,
		NewText("Some more text"))

	out, err := Output(handle, 80)
	if assert.NoError(t, err) {
		// not actually a bug; this is an oddity of Kiselyov's
		// algorithm, that it always breaks outside of a group.
		assert.Equal(t, "Some text\nSome more text", out)
	}
	out, err = Output(handle, 4)
	if assert.NoError(t, err) {
		assert.Equal(t, "Some text\nSome more text", out)
	}

	handle = NewGroup(handle)

	out, err = Output(handle, 80)
	if assert.NoError(t, err) {
		assert.Equal(t, "Some text Some more text", out)
	}
	out, err = Output(handle, 4)
	if assert.NoError(t, err) {
		assert.Equal(t, "Some text\nSome more text", out)
	}
}

func TestCondWithTail(t *testing.T) {
	handle := NewDottedList(NewText("a"),
		NewFuncall("b", NewText("16"), NewText("18")),
		NewText("field"),
		NewFuncall("c", NewText("18")))

	out, err := Output(handle, 80)
	if assert.NoError(t, err) {
		assert.Equal(t, "a.b(16, 18).field.c(18)", out)
	}
	out, err = Output(handle, 4)
	if assert.NoError(t, err) {
		assert.Equal(t, `a.b(16,
    18)
 .field
 .c(18)`, out)
	}
}

func TestInvolved(t *testing.T) {
	handle := NewDottedList(NewFuncall("expr", NewText("5")),
		NewFuncall("add", NewDottedList(NewFuncall("expr", NewText("7")),
			NewFuncall("frob"))),
		NewFuncall("mul", NewDottedList(NewFuncall("expr", NewText("17"))),
			NewFuncall("mul", NewDottedList(NewFuncall("expr", NewText("17")))),
			NewFuncall("mul", NewDottedList(NewFuncall("expr", NewText("17"))))))

	out, err := Output(handle, 180)
	if assert.NoError(t, err) {
		assert.Equal(t, `expr(5).add(expr(7).frob()).mul(expr(17), mul(expr(17)), mul(expr(17)))`, out)
	}
	out, err = Output(handle, 50)
	if assert.NoError(t, err) {
		assert.Equal(t, `expr(5).add(expr(7).frob())
       .mul(expr(17), mul(expr(17)), mul(expr(17)))`, out)
	}
	out, err = Output(handle, 4)
	if assert.NoError(t, err) {
		assert.Equal(t, `expr(5).add(expr(7).frob())
       .mul(expr(17),
            mul(expr(17)),
            mul(expr(17)))`, out)
	}
}
