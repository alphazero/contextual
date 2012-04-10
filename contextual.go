// Package contextual defines the semantics of a generalized
// hierarchical namespace of string names, and untyped values.
package contextual

import (
	"fmt"
)

const (
	/* - general errors - */
	IllegalArgumentError = "ERR - illegal argument"
	NilParentError       = "ERR - parent is nil"
	NilNameError         = "ERR - name is nil/zero-value"
	NegativeNArgError    = "ERR - hierchy walk steps 'n' is negative"

	/* - binding op errors - */
	NilValueError      = "ERR - nil values are not allowed"
	AlreadyBoundError  = "ERR - already bound error"
	NoSuchBindingError = "ERR - no such binding"
)

// General errors
type Error struct {
	msg string
}

func (e Error) Error() string {
	return e.msg
}
func (e Error) Is(errmsg string) bool {
	return e.msg == errmsg
}

// Binding op errors
type BindingError struct {
	err   *Error
	name  string
	value interface{}
}

func newBindingError(msg string, n string, v interface{}) BindingError {
	e := &Error{msg}
	return BindingError{e, n, v}
}
func (e BindingError) Error() string {
	return fmt.Sprintf("%s - (name:%s - value:%v)", e.err.msg, e.name, e.value)
}
func (e BindingError) Is(errmsg string) bool {
	return e.err.msg == errmsg
}

// ----------------------------------------------------------------------------
// Contextual API
// ----------------------------------------------------------------------------

// Contexts are hierarchical namespaces.
// TODO:
//  - count of elements
//  - IsEmpty()
//
type Context interface {
	// Returns true if root context.
	IsRoot() bool

	// Returns true if context is empty.  Note that like Size(), this measure
	// is contextual/relative from the perspective of a child context. (A sibling
	// context may get distinct results.)
	// This method is equivalent to context.Size()==0
	IsEmpty() bool

	// Returns the size of context, which is a count of the visible binding
	// in the context.
	Size() int

	// Lookup will return a non-nil interface{} reference if a non-nil value binding
	// is present in the context or its parental hierarchical path.  The receiver is
	// first checked, and if not root, successive parents (including root) will be searched.
	//
	// Errors:
	//
	//  NilNameError <= zero-value names are not allowed
	Lookup(name string) (value interface{}, e error)

	// LookupN is a constrained variant of Lookup.  (See Lookup() for general details)
	//
	// This method will limit its walk up the hierarchy (if possible) to number of
	// steps (param: n).
	//
	// Errors:
	//
	//  NilNameError <= zero-value names are not allowed
	//  NegativeNArgError <= n is negative
	LookupN(name string, n int) (interface{}, error)

	// Bind will bind the given value to the name in the receiver.
	//
	// Errors:
	//
	//  NilNameError <= zero-value names are not allowed
	//  NilValueError <= nil values are not allowed
	//  AlreadyBoundError <= a value is already bound to the name
	Bind(name string, value interface{}) error

	// Unbind will delete a value binding to the provided name.
	// The unbound value is returned. Unbind is only applicable
	// to bindings in the receiving Context, i.e. any binding
	// that is accessible via LookupN(name, 0).
	//
	// Errors:
	//
	//  NilNameError <= zero-value names are not allowed
	//  NoSuchBindingError <= no values are bound to the name
	Unbind(name string) (unboundValue interface{}, e error)

	// Rebind's semantics are precisely identical to an Unbind followed
	// by a Bound.
	// The unbound value is returned.
	//
	// Errors:
	//
	//  NoSuchBinding <= no values were bound to the name
	//  NilNameError <= zero-value names are not allowed
	//  NilValueError <= nil values are not allowed
	Rebind(name string, value interface{}) (unboundValue interface{}, e error)
}
