// Copyright 2011 Joubin Houshyar.  All rights reserved.
// Use of this source code is governed by a 2-clause BSD
// license that can be found in the LICENSE file.

package contextual

import (
	//	"errors"
	"log" // TEMP
)

type context struct {
	parent   *context
	bindings map[string]interface{}
}

func newContext() *context {
	return &context{nil, make(map[string]interface{})}
}

// NewContext makes and initializes a new root context
func NewContext() *context {
	return newContext()
}

// REVU: hmm .. c.newChild() or this?  (concern is security)
func ChildContext(p *context) (c *context, e Error) {
	if p == nil {
		return nil, error{NilParentError}
	}

	c = newContext()
	c.parent = p
	return
}

func (c *context) IsRoot() bool {
	return c.parent == nil
}

// Returns true if context is empty.  Note that like Count(), this measure
// is contextual/relative from the perspective of a child context. (A sibling
// context may get distinct results.)
func (c *context) IsEmpty() bool {
	if len(c.bindings) > 0 {
		return false
	}
	if c.parent != nil {
		return c.parent.IsEmpty()
	}
	return true
}

func (c *context) Size() int {
	var c0 int
	if c.parent != nil {
		c0 = c.parent.Size()
	}
	return len(c.bindings) + c0
}

// Per spec:
// Lookup will return a non-nil interface{} reference if a non-nil value binding
// is present in the context or its parental hierarchical path.  The receiver is
// first checked, and if not root, successive parents (including root) will be searched.
//
// Errors:
//
//  NilNameError <= nil names are not allowed
func (c *context) Lookup(name string) (value interface{}, e Error) {
	log.Printf("Lookup(%s)\n", name)
	if name == "" {
		return nil, error{NilNameError}
	}

	if value = c.bindings[name]; value == nil {
		log.Printf("Lookup(%s) = %v\n", name, value)
		for n, v := range c.bindings {
			log.Printf("\tdebug [%s] => %v\n", n, v)
		}
		if c.parent != nil {
			return c.parent.Lookup(name)
		}
	}
	return
}

// Per spec:
// LookupN is a constrained variant of Lookup.  (See Lookup() for general details)
//
// This method will limit its walk up the hierarchy (if possible) to number of
// steps (param: n).
//
// Errors:
//
//  NilNameError <= nil names are not allowed
//  NegativeNArgError <= n is negative
func (c *context) LookupN(name string, n int) (value interface{}, e Error) {
	if name == "" {
		return nil, error{NilNameError}
	}
	if n < 0 {
		return nil, error{NegativeNArgError}
	}

	if value = c.bindings[name]; value == nil {
		n--
		if c.parent != nil && n >= 0 {
			return c.parent.LookupN(name, n)
		}
	}
	return
}

// Per spec:
// Bind will bind the given value to the name in the receiver.
//
// Errors:
//
//  NilNameError <= nil names are not allowed
//  NilValueError <= nil values are not allowed
//  AlreadyBoundError <= a value is already bound to the name
func (c *context) Bind(name string, value interface{}) Error {
	if name == "" {
		return error{NilNameError}
	}
	if value == nil {
		return error{NilValueError}
	}

	if v := c.bindings[name]; v != nil {
		return newBindingError(AlreadyBoundError, name, v)
	}

	c.bindings[name] = value
	return nil
}

// Per spec:
// Unbind will delete a value binding to the provided name.
// The unbound value is returned. Unbind is only applicable
// to bindings in the receiving Context, i.e. any binding
// that is accessible via LookupN(name, 0).
//
// Errors:
//
//  NilNameError <= zero-value names are not allowed
//  NoSuchBindingError <= no values are bound to the name
func (c *context) Unbind(name string) (value interface{}, e Error) {
	if name == "" {
		return nil, error{NilValueError}
	}
	if value = c.bindings[name]; value == nil {
		return nil, newBindingError(NoSuchBindingError, name, value)
	}

	delete(c.bindings, name)
	return
}

// Per spec:
// Rebind's semantics are precisely identical to an Unbind followed
// by a Bound.
// The unbound value is returned.
//
// Errors:
//
//  NoSuchBinding <= no values were bound to the name
//  NilNameError <= zero-value names are not allowed
//  NilValueError <= nil values are not allowed
func (c *context) Rebind(name string, value interface{}) (unboundValue interface{}, e Error) {

	if unboundValue, e = c.Unbind(name); e != nil {
		return
	}
	e = c.Bind(name, value)

	return
}
