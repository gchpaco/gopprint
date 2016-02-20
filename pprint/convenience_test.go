package pprint

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCSV(t *testing.T) {
	assert.Equal(t, `Nest(Text("Foo")Text(",")Cond(" ","","")Text("Bar"))`,
		NewCSV(NewText("Foo"), NewText("Bar")).String())
	assert.Equal(t, `Nest(Text("Foo"))`, NewCSV(NewText("Foo")).String())
	assert.Equal(t, `Text("")`, NewCSV().String())
}

func TestArgs(t *testing.T) {
	assert.Equal(t, `Text("(")Nest(Text("Foo")Text(",")Cond(" ","","")Text("Bar"))Text(")")`,
		NewArgs(NewText("Foo"), NewText("Bar")).String())
	assert.Equal(t, `Text("(")Nest(Text("Foo"))Text(")")`, NewArgs(NewText("Foo")).String())
	assert.Equal(t, `Text("(")Text("")Text(")")`, NewArgs().String())
}

func TestDottedList(t *testing.T) {
	assert.Equal(t, `Text("Foo")Nest(Text(".")Text("Bar")Cond(".",".","")Text("Baz"))`,
		NewDottedList(NewText("Foo"), NewText("Bar"), NewText("Baz")).String())
	assert.Equal(t, `Nest(Text("Foo"))`, NewDottedList(NewText("Foo")).String())
	assert.Equal(t, `Text("")`, NewDottedList().String())
}
