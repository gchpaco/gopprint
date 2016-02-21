package pprint

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCSV(t *testing.T) {
	assert.Equal(t, `Nest(Text("Foo")Text(",")Cond(" ","","")Text("Bar"))`,
		CSV(Text("Foo"), Text("Bar")).String())
	assert.Equal(t, `Nest(Text("Foo"))`, CSV(Text("Foo")).String())
	assert.Equal(t, `Text("")`, CSV().String())
}

func TestArgs(t *testing.T) {
	assert.Equal(t, `Text("(")Nest(Text("Foo")Text(",")Cond(" ","","")Text("Bar"))Text(")")`,
		Args(Text("Foo"), Text("Bar")).String())
	assert.Equal(t, `Text("(")Nest(Text("Foo"))Text(")")`, Args(Text("Foo")).String())
	assert.Equal(t, `Text("(")Text("")Text(")")`, Args().String())
}

func TestDottedList(t *testing.T) {
	assert.Equal(t, `Text("Foo")Nest(Text(".")Text("Bar")Cond(".",".","")Text("Baz"))`,
		DottedList(Text("Foo"), Text("Bar"), Text("Baz")).String())
	assert.Equal(t, `Nest(Text("Foo"))`, DottedList(Text("Foo")).String())
	assert.Equal(t, `Text("")`, DottedList().String())
}
