// Copyright 2011 Joubin Houshyar.  All rights reserved.
// Use of this source code is governed by a 2-clause BSD
// license that can be found in the LICENSE file.

package contextual

import (
	"fmt"
)

/*
 * misusing example as a test as the output mechanism makes it close
 * to a simple and convenient specification based testing.
 */

// Demonstrates creation and inspection of Error.
func ExampleError() {

	cerr := error{NilNameError}
	if cerr.Is(NilNameError) {
		fmt.Println(cerr)
	}

	// Output:
	// ERR - name is nil/zero-value
}

// Demonstrates creation and inspection of BindingError.
func ExampleBindingError() {

	name := "woof"
	value := "snowy"

	cerr := newBindingError(NilValueError, name, nil)
	if cerr.Is(NilValueError) {
		fmt.Println(cerr)
	}

	cerr = newBindingError(AlreadyBoundError, name, value)
	if cerr.Is(AlreadyBoundError) {
		fmt.Println(cerr)
	}

	cerr = newBindingError(NoSuchBindingError, name, nil)
	if cerr.Is(NoSuchBindingError) {
		fmt.Println(cerr)
	}

	// Output:
	// ERR - nil values are not allowed - (name:woof - value:<nil>)
	// ERR - already bound error - (name:woof - value:snowy)
	// ERR - no such binding - (name:woof - value:<nil>)
}
