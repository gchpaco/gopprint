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
	handle := Text("Some text")

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
	handle := Concat(Text("Some text"), CondLB,
		Text("Some more text"))

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

	handle = Group(handle)

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
	handle := DottedList(Text("a"),
		Funcall("b", Text("16"), Text("18")),
		Text("field"),
		Funcall("c", Text("18")))

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
	handle := DottedList(Funcall("expr", Text("5")),
		Funcall("add", DottedList(Funcall("expr", Text("7")),
			Funcall("frob"))),
		Funcall("mul", DottedList(Funcall("expr", Text("17"))),
			Funcall("mul", DottedList(Funcall("expr", Text("17")))),
			Funcall("mul", DottedList(Funcall("expr", Text("17"))))))

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
