package pprint

import "fmt"

// Element is a catch all type for the various pretty printer
// primitives.
type Element interface {
	Width() int
	String() string
}

type Text struct {
	text string
}

func (d *Text) Width() int {
	return len(d.text)
}

func (d *Text) String() string {
	return fmt.Sprintf(`Text("%s")`, d.text)
}

func NewText(text string) Element {
	return &Text{text: text}
}

type Cond struct {
	small, continuation, tail string
}

func (d *Cond) Width() int {
	return len(d.small)
}

func (d *Cond) String() string {
	return fmt.Sprintf(`Cond("%s","%s","%s")`, d.small, d.continuation, d.tail)
}

func NewCond(small, cont, tail string) Element {
	return &Cond{small: small, continuation: cont, tail: tail}
}

type LineBreak struct {
}

func (d *LineBreak) Width() int {
	return 0
}

func (d *LineBreak) String() string {
	return "CR"
}

type Concat struct {
	children []Element
}

func (d *Concat) Width() int {
	w := 0
	for _, elt := range d.children {
		w += elt.Width()
	}
	return w
}

func (d *Concat) String() string {
	w := ""
	for _, elt := range d.children {
		w += elt.String()
	}
	return w
}

func NewConcat(elements ...Element) Element {
	return &Concat{children: elements}
}

type Group struct {
	child Element
}

func (d *Group) Width() int {
	return d.child.Width()
}

func (d *Group) String() string {
	return fmt.Sprintf(`Group(%s)`, d.child.String())
}

func NewGroup(element Element) Element {
	return &Group{child: element}
}

type Nest struct {
	child Element
}

func (d *Nest) Width() int {
	return d.child.Width()
}

func (d *Nest) String() string {
	return fmt.Sprintf(`Nest(%s)`, d.child.String())
}

func NewNest(element Element) Element {
	return &Nest{child: element}
}

var (
	Empty  = NewText("")
	CondLB = NewCond(" ", "", "")
	DotLB  = NewCond(".", ".", "")
	LB     = new(LineBreak)
)
