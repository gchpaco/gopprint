package pprint

import "fmt"

// Element is a catch all type for the various pretty printer
// primitives.
type Element interface {
	// Width yields how many characters Element would take on a line
	// without wrapping.
	Width() int
	// String renders the Element in a debug-suitable form.
	String() string
	// private here is to make sure other packages cannot add new
	// types; new types will break the tree renderer (which could be
	// worked around) but also most new types you would want to add
	// require new stream primitives.
	private()
}

type text struct {
	text string
}

func (d *text) Width() int {
	return len(d.text)
}

func (d *text) String() string {
	return fmt.Sprintf(`Text("%s")`, d.text)
}

func (d *text) private() {
}

// NewText constructs an Element for the given text string.
func NewText(payload string) Element {
	return &text{text: payload}
}

type cond struct {
	small, continuation, tail string
}

func (d *cond) Width() int {
	return len(d.small)
}

func (d *cond) String() string {
	return fmt.Sprintf(`Cond("%s","%s","%s")`, d.small, d.continuation, d.tail)
}

func (d *cond) private() {
}

// NewCond constructs an Element that, if there is room, will render
// as `small`; if there is not room, it will render as `tail`, a line
// break, any required indentation, and then `cont`.
func NewCond(small, cont, tail string) Element {
	return &cond{small: small, continuation: cont, tail: tail}
}

type linebreak struct {
}

func (d *linebreak) Width() int {
	return 0
}

func (d *linebreak) String() string {
	return "CR"
}

func (d *linebreak) private() {
}

type concat struct {
	children []Element
}

func (d *concat) Width() int {
	w := 0
	for _, elt := range d.children {
		w += elt.Width()
	}
	return w
}

func (d *concat) String() string {
	w := ""
	for _, elt := range d.children {
		w += elt.String()
	}
	return w
}

func (d *concat) private() {
}

// NewConcat concatenates `elements` into a new Element.
func NewConcat(elements ...Element) Element {
	return &concat{children: elements}
}

type group struct {
	child Element
}

func (d *group) Width() int {
	return d.child.Width()
}

func (d *group) String() string {
	return fmt.Sprintf(`Group(%s)`, d.child.String())
}

func (d *group) private() {
}

// NewGroup wraps `element` in a type that ensures all line break
// decisions will be consistent; either they will all break, or all
// not break.
func NewGroup(element Element) Element {
	return &group{child: element}
}

type nest struct {
	child Element
}

func (d *nest) Width() int {
	return d.child.Width()
}

func (d *nest) String() string {
	return fmt.Sprintf(`Nest(%s)`, d.child.String())
}

func (d *nest) private() {
}

// NewNest wraps `element` in a type similar to `NewGroup` that
// ensures all line break decisions will be consistent, and also
// enforces that any line break must indent at least as much as the
// start of the `NewNest` element.
func NewNest(element Element) Element {
	return &nest{child: element}
}

var (
	// Empty is a convenience variable for an empty Element.
	Empty = NewText("")
	// CondLB is a convenience variable for a conditional space-or-linebreak.
	CondLB = NewCond(" ", "", "")
	// DotLB is a convenience variable for a conditional
	// dot-or-linebreak-and-then-dot; for example for formatting
	// chained methods.
	DotLB = NewCond(".", ".", "")
	// LB is an unconditional line break.
	LB = new(linebreak)
)
