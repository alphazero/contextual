// Copyright 2011 Joubin Houshyar.  All rights reserved.
// Use of this source code is governed by a 2-clause BSD
// license that can be found in the LICENSE file.

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

type Error interface {
	Is(errmsg string) bool
}

type error struct {
	msg string
}

func (e error) Error() string {
	return e.msg
}
func (e error) Is(errmsg string) bool {
	return e.msg == errmsg
}

// Binding op errors
type bindingError struct {
	err   *error
	name  string
	value interface{}
}

func newBindingError(msg string, n string, v interface{}) Error {
	e := &error{msg}
	return bindingError{e, n, v}
}
func (e bindingError) Error() string {
	return fmt.Sprintf("%s - (name:%s - value:%v)", e.err.msg, e.name, e.value)
}
func (e bindingError) Is(errmsg string) bool {
	return e.err.msg == errmsg
}

// ----------------------------------------------------------------------------
// Contextual API
// ----------------------------------------------------------------------------

// Contexts are hierarchical namespaces.
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

	// Returns a non-negative value of the nesting order (depth) of the context.
	// If IsRoot() is true, depth is 0.
	Depth() int

	// Lookup will return a non-nil interface{} reference if a non-nil value binding
	// is present in the context or its parental hierarchical path.  The receiver is
	// first checked, and if not root, successive parents (including root) will be searched.
	//
	// Errors:
	//
	//  NilNameError <= zero-value names are not allowed
	Lookup(name string) (value interface{}, e Error)

	// LookupN is a constrained variant of Lookup.  (See Lookup() for general details)
	//
	// This method will limit its walk up the hierarchy (if possible) to number of
	// steps (param: n).
	//
	// Errors:
	//
	//  NilNameError <= zero-value names are not allowed
	//  NegativeNArgError <= n is negative
	LookupN(name string, n int) (interface{}, Error)

	// Bind will bind the given value to the name in the receiver.
	//
	// Errors:
	//
	//  NilNameError <= zero-value names are not allowed
	//  NilValueError <= nil values are not allowed
	//  AlreadyBoundError <= a value is already bound to the name
	Bind(name string, value interface{}) Error

	// Unbind will delete a value binding to the provided name.
	// The unbound value is returned. Unbind is only applicable
	// to bindings in the receiving Context, i.e. any binding
	// that is accessible via LookupN(name, 0).
	//
	// Errors:
	//
	//  NilNameError <= zero-value names are not allowed
	//  NoSuchBindingError <= no values are bound to the name
	Unbind(name string) (unboundValue interface{}, e Error)

	// Rebind's semantics are precisely identical to an Unbind followed
	// by a Bound.
	// The unbound value is returned.
	//
	// Errors:
	//
	//  NoSuchBinding <= no values were bound to the name
	//  NilNameError <= zero-value names are not allowed
	//  NilValueError <= nil values are not allowed
	Rebind(name string, value interface{}) (unboundValue interface{}, e Error)
}

// Contextual things
type Contextual interface {
	SetContext(ctx Context)
}

// A component is contextual
// TOOD: ports {in, out}
type Component interface {
	Contextual
}

// A component that is a containment context
type Container interface {
	Component

	Add(c Component) Error
	Remove(c Component) Error
}
