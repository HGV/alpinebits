package alpinebits

import (
	"fmt"
	"strconv"
)

// Path locates an element of a message in error messages, e.g.
// "RatePlan[0].BookingRule[2]". The empty Path means the whole message.
type Path string

// At returns the path of the i-th name element: At("RatePlan", 0) is
// "RatePlan[0]".
func At(name string, i int) Path {
	return Path(name + "[" + strconv.Itoa(i) + "]")
}

// At returns the path of the i-th name element inside p.
func (p Path) At(name string, i int) Path {
	return p.join(At(name, i))
}

// Field returns the path of the name element inside p.
func (p Path) Field(name string) Path {
	return p.join(Path(name))
}

func (p Path) join(child Path) Path {
	if p == "" {
		return child
	}
	return p + "." + child
}

// Addf appends an application error (Type 13) about the element at path at.
func (e *Errors) Addf(code ErrCode, at Path, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	if at != "" {
		msg = string(at) + ": " + msg
	}
	*e = append(*e, Error{Type: ErrTypeApplicationError, Code: code, Message: msg})
}
