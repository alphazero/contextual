package contextual

import (
	"errors"
)

type context struct {
	parent  *context
	objects map[string]interface{}
}

func newContext() *context {
	return &context{nil, make(map[string]interface{})}
}

func RootContext() *context {
	return newContext()
}

func ChildContext(p *context) (c *context, e error) {
	if p == nil {
		return nil, errors.New("")
	}
	c = newContext()
	c.parent = p
	return
}

func (c *context) IsRoot() bool {
	return c.parent == nil
}
