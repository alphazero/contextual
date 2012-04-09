/* Note on file organization:  tests are invoked in sequence by Go test
 * harness and that fact is used here to affirm a set of assumptions as
 * we test the constructs.  As assumptions are affirmed, they are noted
 * in comments below.  So, do not change the order of the tests.
 */

package contextual

import (
	"fmt"
	"testing"
)

// ============================================================================
// testing: contextual.context
// ============================================================================

// NOP - just feedback for test runs. std. per each construct
func TestContextStructStart_NOP(t *testing.T) {
	fmt.Println("contextual.context")
}

// helper
type emptyStruct struct{}

// helper
func mixedTypeValueSet() []interface{} {
	num := 10
	str := emptyStruct{}
	ptr := &emptyStruct{}
	ch := make(chan emptyStruct)
	fn := func() {} // func() is "uncomparable type" oh, well.
	pfn := &fn
	txt := "hello there"

	return []interface{}{
		num, str, ptr, ch /*fn,*/, pfn, txt,
	}

}

// helper
func genericUniqueIndexNames(n int) []string {
	var names []string = make([]string, n)
	var sanitycheck = make(map[string]int)

	for i := 0; i < n; i++ {
		names[i] = fmt.Sprintf("value[%d]", i)
		sanitycheck[names[i]] = i
	}
	// make sure names are indeed unique
	if len(sanitycheck) != len(names) {
		panic("genericUniqueIndexNames")
	}

	return names
}

/* --- CHECK a1 ---------------------------------------------------------------
 *  a1: correct construction and intial state
 *
 * tests for correct construction and initialization of contexts, and asserts
 * on input params of the associated functions.
 * - Context#Size()    // on init
 * - Context#IsEmpty() // on init
 *
 * assumptions:
 * - none
 */

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

func TestNewContextInit(t *testing.T) {
	ctx := NewContext()
	if ctx.IsEmpty() != true {
		t.Fatalf("New Context#IsEmpty returned false")
	}
	if ctx.Size() != 0 {
		t.Fatalf("New Context#Size returned non-zero")
	}

	fmt.Println("\ttest context init - root")
}

func TestNewContextChildInit(t *testing.T) {
	rootCtx := NewContext()

	ctx, _ := ChildContext(rootCtx)
	if ctx.IsEmpty() != true {
		t.Fatalf("New child Context#IsEmpty returned false")
	}
	if ctx.Size() != 0 {
		t.Fatalf("New child Context#Size returned non-zero")
	}

	fmt.Println("\ttest context init - child")
}

/* --- CONFIRMED a1 ----------------------------------------------------------*/

/* --- CHECK a2 ---------------------------------------------------------------
 * - a2: correct parent/child order and IsRoot
 *
 * tests for parent child relationship and correct behavior for
 * Context#IsRoot()
 *
 * assumptions:
 * - a1
 */

// tests contextual.context's compliance with Context#IsRoot()
func TestIsRoot(t *testing.T) {
	// a1 - don't bother with error/nil checks
	rootCtx := NewContext()
	if rootCtx.IsRoot() != true {
		t.Fatalf("IsRoot() for a root context must return true")
	}
	if ctx, _ := ChildContext(rootCtx); ctx.IsRoot() {
		t.Fatalf("IsRoot() for a child context must return false")
	}

	fmt.Println("\tContext#IsRoot()")
}

/* --- CONFIRMED a2 ----------------------------------------------------------*/

/* --- CHECK a3 ---------------------------------------------------------------
 * - a3: correct basic ops for a single (root) context
 *
 * tests for correct behavior of the following Context methods:
 * - Context#Bind()
 * - Context#Lookup()
 * - Context#Size()    // post init
 * - Context#IsEmpty() // post init
 * - Context#Unbind()
 * - Context#Rebind()
 * assumptions:
 * - a1
 * - a2
 */

// Confirm spec compliance for faults
// {Lookup, LookupN, Bind, Unbind}
func TestContextSpecdError(t *testing.T) {
	// setup
	ctx := NewContext()

	// Lookup()

	// test specified errors for Lookup
	//  NilNameError <= zero-value names are not allowed
	if _, e := ctx.Lookup(""); e == nil {
		t.Fatalf("Lookup(nil) expected error: %s", NilNameError)
	}
	v, e := ctx.Lookup("no-such-binding")
	if e != nil {
		t.Fatalf("Unexpected error: %s", e)
	}
	if v != nil {
		t.Fatalf("Lookup(\"\") - expected:%v got:%v", nil, v)
	}

	// LookupN()

	//  NilNameError <= nil names are not allowed
	//  NegativeNArgError <= n is negative
	if _, e := ctx.LookupN("", 0); e == nil {
		t.Fatalf("Lookup(nil) expected error: %s", NilNameError)
	}
	v, e = ctx.LookupN("no-such-binding", 0)
	if e != nil {
		t.Fatalf("Unexpected error: %s", e)
	}
	if v != nil {
		t.Fatalf("Lookup(\"\") - expected:%v got:%v", nil, v)
	}

	// Bind()

	// test specified errors for Bind
	//  NilNameError <= zero-value names are not allowed
	//  NilValueError <= nil values are not allowed
	//  AlreadyBoundError <= a value is already bound to the name
	if e := ctx.Bind("", "some value"); e == nil {
		t.Fatalf("Bind(\"\") expected error: %s", NilNameError)
	}
	if e := ctx.Bind("some key", nil); e == nil {
		t.Fatalf("Bind(\"\") expected error: %s", NilValueError)
	}

	// Unbind()

	// test specified errors for Unbind
	//  NilNameError <= zero-value names are not allowed
	//  NoSuchBindingError <= no values are bound to the name
	wat, e := ctx.Unbind("")
	if e == nil {
		t.Fatalf("Unbind(\"\") expected error: %s", NilNameError)
	}
	if wat != nil {
		t.Fatalf("Unexpected value on faulted return: %s", wat)
	}
	wat, e = ctx.Unbind("some key")
	if e == nil {
		t.Fatalf("Unbind(\"\") expected error: %s", NoSuchBindingError)
	}
	if wat != nil {
		t.Fatalf("Unexpected value on faulted return: %s", wat)
	}

	// Rebind()

	// test specified errors for Rebind
	//  NoSuchBinding <= no values were bound to the name
	//  NilNameError <= zero-value names are not allowed
	//  NilValueError <= nil values are not allowed
	wat, e = ctx.Rebind("", "some value")
	if e == nil {
		t.Fatalf("Rebind(\"\", v) expected error: %s", NilNameError)
	}
	if wat != nil {
		t.Fatalf("Unexpected value on faulted return: %s", wat)
	}
	wat, e = ctx.Rebind("some key", "some value")
	if e == nil {
		t.Fatalf("Rebind(\"\", v) expected error: %s", NoSuchBindingError)
	}
	if wat != nil {
		t.Fatalf("Unexpected value on faulted return: %s", wat)
	}
	if wat, e = ctx.Rebind("doesn't matter", nil); e == nil {
		t.Fatalf("Rebind (nil) expected error: %s", NilValueError)
	}
}

func TestSingleContext(t *testing.T) {
	// setup
	ctx := NewContext()
	values := mixedTypeValueSet()
	names := genericUniqueIndexNames(len(values))

	// Bind()

	for i, name := range names {
		value := values[i]
		if e := ctx.Bind(name, value); e != nil {
			t.Fatalf("Unexpected error: %s", e)
		}
	}

	// Lookup()

	for i, name := range names {
		expv := values[i]
		v, e := ctx.Lookup(name)
		if e != nil {
			t.Fatalf("Unexpected error: %s", e)
		}
		if v != expv {
			t.Fatalf("Lookup(%s) - expected:%s got:%s", name, expv, v)
		}
	}

	// IsEmpty()

	if b := ctx.IsEmpty(); b {
		t.Fatalf("IsEmpty() - expected:%s got:%s", false, b)
	}

	// Size()

	if n := ctx.Size(); n != len(names) {
		t.Fatalf("Size() - expected:%s got:%s", len(names), n)
	}

	// Unbind()

	for i, name := range names {
		expv := values[i]
		v, e := ctx.Unbind(name)
		if e != nil {
			t.Fatalf("Unexpected error: %s", e)
		}
		if v != expv {
			t.Fatalf("Unbind(%s) - expected:%s got:%s", name, expv, v)
		}
	}
	// confirm Size and IsEmpty
	if !ctx.IsEmpty() {
		t.Fatalf("IsEmpty() returned true")
	}
	if ctx.Size() != 0 {
		t.Fatalf("Size() returned non-zero")
	}

	// Rebind()
	ctx.Bind(names[0], values[0]) // setup
	oldv, e := ctx.Rebind(names[0], "on the rebound")
	if e != nil {
		t.Fatalf("unexpected error: %s", e)
	}
	if oldv != values[0] {
		t.Fatalf("Rebind() old value - expected:%v got:%v", values[0], oldv)
	}

	fmt.Println("\tContext compliance - root")
}

/* --- CONFIRMED a3 ----------------------------------------------------------*/

/* --- CHECK a4 ---------------------------------------------------------------
 * - a4: correct basic ops for context hierarchy
 *
 * tests for correct behavior of the following Context methods:
 * - Context#Bind()
 * - Context#Lookup()
 * - Context#Size()    // post init
 * - Context#IsEmpty() // post init
 * - Context#Unbind()
 * - Context#Rebind()
 * assumptions:
 * - a1
 * - a2
 * - a3
 */

func TestContextHierarchy(t *testing.T) {
	// setup
	//	ctx := NewContext()
	//	child1, _ := ChildContext(ctx)
	//	child1_1, _ := ChildContext(child1)
	//	child2, _ := ChildContext(ctx)
	//
	//	values := mixedTypeValueSet()
	//	names := genericUniqueIndexNames(len(values))
}

/* --- CONFIRMED a4 ----------------------------------------------------------*/
