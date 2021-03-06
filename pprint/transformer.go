package pprint

import (
	"io"
	"strings"
)

// toStream recursively converts a document into the stream elements
// we'll be using.  We use channels to organize the coroutines.
func toStream(document Element) <-chan streamElt {
	ch := make(chan streamElt)
	go func() {
		defer close(ch)
		visitElement(document, ch)
	}()
	return ch
}

func visitElement(document Element, out chan<- streamElt) {
	switch doc := document.(type) {
	case *text:
		out <- &textElt{elt{-1}, doc.text}
	case *cond:
		out <- &condElt{elt{-1}, doc.small, doc.continuation, doc.tail}
	case *linebreak:
		out <- &crlfElt{elt{-1}}
	case *concat:
		for _, elt := range doc.children {
			visitElement(elt, out)
		}
	case *group:
		out <- &gbegElt{elt{-1}}
		visitElement(doc.child, out)
		out <- &gendElt{elt{-1}}
	case *nest:
		out <- &nbegElt{elt{-1}}
		out <- &gbegElt{elt{-1}}
		visitElement(doc.child, out)
		out <- &gendElt{elt{-1}}
		out <- &nendElt{elt{-1}}
	default:
		panic("Couldn't understand document type")
	}
}

// annotateLastChar is the next step; it takes the stream elements
// from `toStream` and adds information about the horizontal position
// of their last character.  This is not possible with NBeg and GBeg
// elements as we haven't got enough information yet.
func annotateLastChar(in <-chan streamElt) <-chan streamElt {
	ch := make(chan streamElt)
	go func() {
		defer close(ch)
		position := 0
		for {
			select {
			case elt, ok := <-in:
				if !ok {
					return
				}
				switch elt := elt.(type) {
				case *textElt:
					position += len(elt.payload)
					elt.hpos = position
					ch <- elt
				case *condElt:
					position += len(elt.small)
					elt.hpos = position
					ch <- elt
				case *crlfElt:
					elt.hpos = position
					ch <- elt
				case *gbegElt, *nbegElt:
					// Don't have enough information yet to do this
					// accurately.
					ch <- elt
				case *gendElt:
					elt.hpos = position
					ch <- elt
				case *nendElt:
					elt.hpos = position
					ch <- elt
				}
			}
		}
	}()
	return ch
}

type lookaheads [][]streamElt

func (s lookaheads) addToLast(element ...streamElt) {
	last := len(s) - 1
	s[last] = append(s[last], element...)
}

func (s lookaheads) pushNew() lookaheads {
	return append(s, make([]streamElt, 0))
}

func (s lookaheads) pop() ([]streamElt, lookaheads) {
	last := len(s) - 1
	return s[last], s[0:last]
}

// annotateGBeg is the next step; we take the horizontal position
// information gotten from `annotateLastChar` and compute the `hpos`
// for GBeg elements.  We don't need to do it for NBeg, but for GBeg
// it matters for linebreaks.
func annotateGBeg(in <-chan streamElt) <-chan streamElt {
	ch := make(chan streamElt)
	go func() {
		defer close(ch)
		var lookahead lookaheads
		for {
			select {
			case element, ok := <-in:
				if !ok {
					return
				}
				switch element := element.(type) {
				case *textElt, *condElt, *crlfElt, *nbegElt, *nendElt:
					if len(lookahead) == 0 {
						ch <- element
					} else {
						lookahead.addToLast(element)
					}
				case *gbegElt:
					lookahead = lookahead.pushNew()
				case *gendElt:
					var top []streamElt
					top, lookahead = lookahead.pop()
					if len(lookahead) == 0 {
						// this, then, was the topmost stack
						ch <- &gbegElt{elt{element.hpos}}
						for _, e := range top {
							ch <- e
						}
						ch <- element
					} else {
						lookahead.addToLast(&gbegElt{elt{element.hpos}})
						lookahead.addToLast(top...)
						lookahead.addToLast(element)
					}
				}
			}
		}
	}()
	return ch
}

// Kiselyov's original formulation includes an alternate third phase
// which limits lookahead to the width of the page.  This is difficult
// for us because we don't guarantee docs are of nonzero length,
// although that could be finessed, and also it adds extra complexity
// for minimal benefit.  This implementation skips it.

// The final phase is to compute output.  Each time we see a
// `gbeg_element_t`, we can compare its `hpos` with `rightEdge` to see
// whether it'll fit without breaking.  If it does fit, increment
// `fittingElements` and proceed, which will cause the logic for
// `text_element_t` and `cond_element_t` to just append stuff without
// line breaks.  If it doesn't fit, set `fittingElements` to 0, which
// will cause `cond_element_t` to do line breaks.  When we do a line
// break, we need to compute where the new right edge of the 'page'
// would be in the context of the original stream; so if we saw a
// `cond_element_t` with `e.hpos` of 300 (meaning it ends at
// horizontal position 300), the new right edge would be 300 -
// indentation + page width.
func output(in <-chan streamElt, width int, output io.Writer) error {
	fittingElements := 0
	rightEdge := width
	hpos := 0
	var indent []int
	for {
		select {
		case elt, ok := <-in:
			if !ok {
				return nil
			}
			switch elt := elt.(type) {
			case *textElt:
				_, err := io.WriteString(output, elt.payload)
				if err != nil {
					return err
				}
				hpos += len(elt.payload)
			case *condElt:
				if fittingElements == 0 {
					var currentIndent int
					if len(indent) == 0 {
						currentIndent = 0
					} else {
						currentIndent = indent[len(indent)-1]
					}
					_, err := io.WriteString(output, elt.tail)
					if err != nil {
						return err
					}
					_, err = io.WriteString(output, "\n")
					if err != nil {
						return err
					}
					_, err = io.WriteString(output, strings.Repeat(" ", currentIndent))
					if err != nil {
						return err
					}
					_, err = io.WriteString(output, elt.cont)
					if err != nil {
						return err
					}
					fittingElements = 0
					hpos = currentIndent + len(elt.cont)
					rightEdge = (width - hpos) + elt.hpos
				} else {
					_, err := io.WriteString(output, elt.small)
					if err != nil {
						return err
					}
					hpos += len(elt.small)
				}
			case *crlfElt:
				var currentIndent int
				if len(indent) == 0 {
					currentIndent = 0
				} else {
					currentIndent = indent[len(indent)-1]
				}
				_, err := io.WriteString(output, "\n")
				if err != nil {
					return err
				}
				_, err = io.WriteString(output, strings.Repeat(" ", currentIndent))
				if err != nil {
					return err
				}
				fittingElements = 0
				hpos = currentIndent
				rightEdge = (width - hpos) + elt.hpos
			case *gbegElt:
				if fittingElements != 0 || elt.hpos <= rightEdge {
					fittingElements++
				} else {
					fittingElements = 0
				}
			case *gendElt:
				if fittingElements != 0 {
					fittingElements--
				}
			case *nbegElt:
				indent = append(indent, hpos)
			case *nendElt:
				if len(indent) > 0 {
					indent = indent[0 : len(indent)-1]
				}
			}
		}
	}
}

// PrettyPrint prints `doc` to `out` assuming a right page edge of
// `width`.
func PrettyPrint(doc Element, width int, out io.Writer) error {
	return output(annotateGBeg(annotateLastChar(toStream(doc))), width, out)
}
