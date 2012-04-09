package contextual

import (
	"errors"
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
func ChildContext(p *context) (c *context, e error) {
	if p == nil {
		return nil, errors.New("p is nil")
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
		return true
	}
	if c.parent != nil {
		return c.parent.IsEmpty()
	}
	return false
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
func (c *context) Lookup(name string) (value interface{}, e error) {
	if name == "" {
		return nil, Error{NilNameError}
	}

	if value = c.bindings[name]; value == nil {
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
func (c *context) LookupN(name string, n int) (value interface{}, e error) {
	if name == "" {
		return nil, Error{NilNameError}
	}
	if n < 0 {
		return nil, Error{IllegalArgumentError}
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
func (c *context) Bind(name string, value interface{}) error {
	if name == "" {
		return Error{NilValueError}
	}
	if value == nil {
		return Error{NilNameError}
	}

	if value = c.bindings[name]; value != nil {
		return newBindingError(AlreadyBoundError, name, value)
	}

	c.bindings[name] = value
	return nil
}

// Per spec:
// Unbind will delete a value binding to the provided name.
// The unbound value is returned.
//
// Errors:
//
//  NilNameError <= zero-value names are not allowed
//  NoSuchBindingError <= no values are bound to the name
func (c *context) Unbind(name string) (value interface{}, e error) {
	if name == "" {
		return nil, Error{NilValueError}
	}
	if value = c.bindings[name]; value == nil {
		return nil, newBindingError(NoSuchBindingError, name, value)
	}

	c.bindings[name] = nil

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
func (c *context) Rebind(name string, value interface{}) (unboundValue interface{}, e error) {

	if unboundValue, e = c.Unbind(name); e != nil {
		return
	}
	e = c.Bind(name, value)

	return
}
