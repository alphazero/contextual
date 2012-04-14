// Copyright 2011 Joubin Houshyar.  All rights reserved.
// Use of this source code is governed by a 2-clause BSD
// license that can be found in the LICENSE file.

package contextual

type component struct {
	context Context
}

func NewComponent() Component {
	return &component{}
}
func (c *component) SetContext(ctx Context) {
	c.context = ctx
}
