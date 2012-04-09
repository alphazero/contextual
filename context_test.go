package contextual

import (
	"testing"
	"fmt"
)

func TestContextConstruct(t *testing.T) {
	var ctx Context = RootContext()
	if ctx == nil {
		t.Fatalf("RootContext returned nil ref")
	}

	fmt.Println("create root context")
}

func TestChildContextConstruct(t *testing.T) {
	rootCtx := RootContext()

	ctx, e := ChildContext(rootCtx)
	if e != nil {
		t.Fatalf("Unexpected error: %s", e)
	}
	if ctx == nil {
		t.Fatalf("ChildContext returned nil ref")
	}
	ctx, e = ChildContext(nil)
	if e == nil {
		t.Fatalf("Expecting error on ChildContext(nil)")
	}

	fmt.Println("create child context")
}

// ASSUMPTIONS:
//  a1 - previous passing tests affirming correct constructor behavior
func TestIsRoot(t *testing.T) {
	// a1 - don't bother with error/nil checks
	rootCtx := RootContext(); if rootCtx.IsRoot() != true {
		t.Fatalf("IsRoot() for a root context must return true")
	}
	if ctx, _ := ChildContext(rootCtx); ctx.IsRoot() {
		t.Fatalf("IsRoot() for a child context must return false")
	}

	fmt.Println("check IsRoot()")

}
