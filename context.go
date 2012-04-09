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

func RootContext() *context {
	return newContext()
}

func ChildContext(p *context) (c *context, e error) {
	if p == nil {
		return nil, errors.New("p is nil")
	}
	c = newContext()
	c.parent = p
	return
}

func ___TEMP___(c *context) Context {
	assert(c)
	return c
}
func assert(c *context) {
	if c == nil {
		panic("c is receiver")
	}
}
func (c *context) IsRoot() bool {
	assert(c)
	return c.parent == nil
}

// Lookup will return a non-nil interface{} reference if a non-nil value binding
// is present in the context or its parental hierarchical path.  The receiver is
// first checked, and if not root, successive parents (including root) will be searched.
//
// Errors:
//
//  NilNameError <= nil names are not allowed
func (c *context) Lookup(name string) (value interface{}, e error) {
	assert(c)
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

// LookupN is a constrained variant of Lookup.  (See Lookup() for general details)
//
// This method will limit its walk up the hierarchy (if possible) to number of
// steps (param: n).
//
// Errors:
//
//  NilNameError <= nil names are not allowed
//  IllegalArgument <= n is negative
func (c *context) LookupN(name string, n int) (interface{}, error) {
	assert(c)
	return nil, nil
}

// Bind will bind the given value to the name in the receiver.
//
// Errors:
//
//  NilNameError <= nil names are not allowed
//  NilValueError <= nil values are not allowed
//  AlreadyBoundError <= a value is already bound to the name
func (c *context) Bind(name string, value interface{}) error {
	assert(c)
	return nil
}

// Unbind will delete an value binding to the provided name.
//
// Errors:
//
//  NilNameError <= nil names are not allowed
//  NoSuchBinding <= no values are bound to the name
func (c *context) Unbind(name string) error {
	assert(c)
	return nil
}

// Rebind's semantics are precisely identical to an Unbind followed
// by a Bound.
//
// Errors:
//
//  NoSuchBinding <= no values were bound to the name
//  NilNameError <= nil names are not allowed
//  NilValueError <= nil values are not allowed
func (c *context) Rebind(name string, value interface{}) error {
	assert(c)
	return nil
}
