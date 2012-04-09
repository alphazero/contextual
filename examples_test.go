package contextual

import (
	"fmt"
)

/*
 * misusing example as a test as the output mechanism makes it close
 * to a simple and convenient specification based testing.
 */

// demonstrate both creation and inspection of errors
func ExampleBindingError() {

	name := "woof"
	value := "snowy"
	var cerr BindingError

	cerr = newBindingError(NilValueError, name, nil)
	fmt.Println(cerr.Error())
	fmt.Println(cerr.Is(NilValueError))

	cerr = newBindingError(AlreadyBoundError, name, value)
	fmt.Println(cerr.Error())
	fmt.Println(cerr.Is(AlreadyBoundError))

	cerr = newBindingError(NoSuchBindingError, name, nil)
	fmt.Println(cerr.Error())
	fmt.Println(cerr.Is(NoSuchBindingError))

	// Output:
	// ERR - nil values are not allowed - (name:woof - value:<nil>)
	// true
	// ERR - already bound error - (name:woof - value:snowy)
	// true
	// ERR - no such binding - (name:woof - value:<nil>)
	// true
}
