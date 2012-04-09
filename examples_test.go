package contextual

import (
	"fmt"
)

func ExampleBindingError() {

	name := "woof"
	value := "snowy"
	var cerr BindingError

	cerr = BindingError{name, nil, NilValueError}
	fmt.Println(cerr.Error())

	cerr = BindingError{name, value, AlreadyBoundError}
	fmt.Println(cerr.Error())

	// Output:
	// ERR - nil values are not allowed - (name:woof - value:<nil>)
	// ERR - already bound error - (name:woof - value:snowy)
}
