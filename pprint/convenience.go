package pprint

var (
	comma  = NewText(",")
	dot    = NewText(".")
	lparen = NewText("(")
	rparen = NewText(")")
)

// NewCSV wraps `elements` with a comma separated list.
func NewCSV(elements ...Element) Element {
	if len(elements) == 0 {
		return Empty
	}
	elts := make([]Element, len(elements)*3-2)
	pos := 0
	for _, elt := range elements {
		if pos == 0 {
			elts[pos] = elt
			pos++
		} else {
			elts[pos] = comma
			elts[pos+1] = CondLB
			elts[pos+2] = elt
			pos += 3
		}
	}
	return NewNest(NewConcat(elts...))
}

// NewArgs formats `elements` in a manner suitable for C style
// arguments.
func NewArgs(elements ...Element) Element {
	return NewConcat(lparen, NewCSV(elements...), rparen)
}

// NewDottedList formats `elements` in a manner suitable for chained
// method calls, á la "fluent" interfaces.
func NewDottedList(elements ...Element) Element {
	if len(elements) == 0 {
		return Empty
	} else if len(elements) == 1 {
		return NewNest(elements[0])
	}
	elts := make([]Element, len(elements)*2-1)
	pos := 0
	for _, elt := range elements {
		if pos == 0 {
			elts[pos] = elt
			pos++
		} else if pos == 1 {
			// we don't want to break on the first dot; it's ugly.
			elts[pos] = dot
			elts[pos+1] = elt
			pos += 2
		} else {
			elts[pos] = DotLB
			elts[pos+1] = elt
			pos += 2
		}
	}
	// Bit involved; we want NewDottedList(a, b, c) to turn into
	// Concat(a, Nest(Concat(".", b, DotLB, c))).  We want this
	// because it means that the dots line up on linebreaks nicely.
	return NewConcat(elts[0], NewNest(NewConcat(elts[1:]...)))
}

// NewFuncall formats `args` as a function call of the function
// `name`.
func NewFuncall(name string, args ...Element) Element {
	return NewConcat(NewText(name), NewArgs(args...))
}
