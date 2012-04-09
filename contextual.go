// Package contextual defines the semantics of a generalized
// hierarchical namespace of string names, and untyped values.
package contextual

import (
	"fmt"
)

// REVU:
// don't like this - need a more general approach.
const (
	/* - general errors - */
	IllegalArgumentError = "ERR - illegal argument"
	NilParentError       = "ERR - parent is nil"
	NilNameError         = "ERR - name is nil/zero-value"
	NegativeNArrError    = "ERR - hierchy walk steps 'n' is negative"
	NoSuchBindingError   = "ERR - no such binding"

	/* - binding op errors - */
	NilValueError     = "ERR - nil values are not allowed"
	AlreadyBoundError = "ERR - already bound error"
)

// General errors
type Error struct {
	msg string
}

func (e Error) Error() string {
	return e.msg
}

// Binding op errors
type BindingError struct {
	name  string
	value interface{}
	msg   string
}

func (e *BindingError) Error() string {
	return fmt.Sprintf("%s - (name:%s - value:%v)", e.msg, e.name, e.value)
}

// REVU-END

// Contexts are hierarchical namespaces.
type Context interface {
	// Returns true if root context.
	IsRoot() bool

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

	// Unbind will delete an value binding to the provided name.
	//
	// Errors:
	//
	//  NilNameError <= zero-value names are not allowed
	//  NoSuchBinding <= no values are bound to the name
	Unbind(name string) error

	// Rebind's semantics are precisely identical to an Unbind followed
	// by a Bound.
	//
	// Errors:
	//
	//  NoSuchBinding <= no values were bound to the name
	//  NilNameError <= zero-value names are not allowed
	//  NilValueError <= nil values are not allowed
	Rebind(name string, value interface{}) error
}
