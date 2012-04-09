package contextual

import (
	"fmt"
	"testing"
)

func TestContextStructStart_NOP(t *testing.T) {
	fmt.Println("contextual.context")
}
func TestNewContext(t *testing.T) {
	var ctx Context = NewContext()
	if ctx == nil {
		t.Fatalf("NewContext returned nil ref")
	}

	fmt.Println("\tcreate root context")
}

func TestNewContextChild(t *testing.T) {
	rootCtx := NewContext()

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

	fmt.Println("\tcreate child context")
}

// ASSUMPTIONS:
//  a1 - previous passing tests affirming correct constructor behavior

func TestIsRoot(t *testing.T) {
	// a1 - don't bother with error/nil checks
	rootCtx := NewContext()
	if rootCtx.IsRoot() != true {
		t.Fatalf("IsRoot() for a root context must return true")
	}
	if ctx, _ := ChildContext(rootCtx); ctx.IsRoot() {
		t.Fatalf("IsRoot() for a child context must return false")
	}

	fmt.Println("\tcheck IsRoot()")
}
