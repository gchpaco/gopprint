package pprint

import (
	"fmt"
)

// While the document types are handy for creating a layout, they're
// not that useful for actually pretty printing the document.  For
// that, we use these stream types.

type streamElt interface {
	String() string
}

type elt struct {
	hpos int
}

type textElt struct {
	elt
	payload string
}

func (e *textElt) String() string {
	return fmt.Sprintf(`TE(%d,"%s")`, e.hpos, e.payload)
}

type condElt struct {
	elt
	small, cont, tail string
}

func (e *condElt) String() string {
	return fmt.Sprintf(`CE(%d,"%s","%s","%s")`, e.hpos, e.small, e.cont, e.tail)
}

type crlfElt struct {
	elt
}

func (e *crlfElt) String() string {
	return fmt.Sprintf(`CR(%d)`, e.hpos)
}

type nbegElt struct {
	elt
}

func (e *nbegElt) String() string {
	return fmt.Sprintf(`NBeg(%d)`, e.hpos)
}

type nendElt struct {
	elt
}

func (e *nendElt) String() string {
	return fmt.Sprintf(`NEnd(%d)`, e.hpos)
}

type gbegElt struct {
	elt
}

func (e *gbegElt) String() string {
	return fmt.Sprintf(`GBeg(%d)`, e.hpos)
}

type gendElt struct {
	elt
}

func (e *gendElt) String() string {
	return fmt.Sprintf(`GEnd(%d)`, e.hpos)
}
