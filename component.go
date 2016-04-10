// Copyright 2011-2016 Joubin Houshyar.  All rights reserved.
// Use of this source code is governed by a 2-clause BSD
// license that can be found in the LICENSE file.

package contextual

// a component is a Contextual object
type component struct {
	context Context
}

func NewComponent() Component {
	return &component{}
}
func (c *component) SetContext(ctx Context) {
	c.context = ctx
}
