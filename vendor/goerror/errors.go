// Copyright 2010-2016 Joubin Houshyar.  All rights reserved.
// Use of this source code is governed by a 2-clause BSD
// license that can be found in the LICENSE file.

// package goerror provides for creation of semantic error types.
// Using the stdlib error package, when trying to discern the error
// type returned by a function, we can either (a) compare to a well-known
// error reference (for example, io.EOF), or (b) we have to parse the
// actual error message. The former, e.g. io.EOF, is fine for basic cases,
// but obviously won't allow for call-site specific information in the error.
//
// For example, let's say we create an error for asserting input arg
// correctness, something like an AsssertError, or IllegalArgumentError. We
// can try pattern (a) per io.EOF, in which case we can certainly return that
// error, but can't provide additional info such as which precise arg caused
// the error. Or we can return a plain jane error with a formatted message,
// in which case we can't immediately tell what 'kind' of error was returned.
//
// This package addresses this concern by providing error 'types' that can
// be generically defined at (some) package level and then used with explicit
// additional details.
//
// Errors are created using 'Define'. (Note, not 'New', since this merely
// defines an error type).
//
//     var (
//         TerribleError      = goerror.Define("TerribleError")
//         NotSoTerribleError = goerror.Define("NotSoTerribleError")
//     )
//
// Such error types can then be 'instantiated' using the defintion,
// wherever one would normally create and/or return a generic error.
//
//     // function foo may return either TerribleError or
//     // NotSoTerribleError
//     func foo() error {
//         // ...
//         if flipcoin() {
//            return TerribleError("an example usage")
//         }
//         return NotSoTerribleError() // detailed info is optional
//     }
//
// And in the functional callsite, we can specifically check to see
// what type of error we got.
//
//    if e := foo(); e != nil {
//        switch typ := goerror.TypeOf(e); {
//        case typ.Is(TerribleError):
//            /* handle it */
//        case typ.Is(NotSoTerribleError):
//            /* handle it */
//        }
//    }
//
package goerror

import (
	"errors"
	"fmt"
)

// ideally we want optional args capped at 1 item
// but you can't do that in go. So the error generating
// functions will create an error string that is simply
// a concatenation of the args passed here.
type errFn func(...string) *Error

// Note that this type is exported *only* in order to surface Is() to package docs.
// Otherwise, package users should not directly use this type.
type Error struct {
	error
	cause error
}

// defines a new categorical error.
func Define(category string) errFn {
	return func(args ...string) *Error {
		errstr := category
		if len(args) == 0 {
			goto done
		}
		errstr += " - "
		// since nargs can be > 1 might as well
		// make a virtue of it and pretty concat the args
		for _, arg := range args {
			errstr += arg
			errstr += " "
		}
		errstr = errstr[:len(errstr)-1]
	done:
		return &Error{error: errors.New(errstr)}
	}
}

// Returns an Error, typically for use in conjunction
// with the Error#Is(). Function name is as such to
// allow for a readable call site, as below:
//
//     if goerror.TypeOf(e).Is(AssertionError)
//
// If the input arg 'e' is a plain (builtin) error, it is
// converted to a goerror.Error pointer.
func TypeOf(e error) *Error {
	if e0, ok := e.(*Error); ok {
		return e0
	}
	return &Error{error: e}
}

// Returns associated cause, or nil.
func (e *Error) Cause() error {
	return e.cause
}

// Associate a root cause error with the given error.
// If cause is already set, subsequent calls to this function
// are ignored.
//
// Typcial usage pattern:
//
//
//   var WriteError = goerror.Define("Write Error")
//
//   func writeBuffer(..) error {
//      ...
//      // let's pretend we did io and got an error 'ioerror'
//      return WriteError("in writeBuffer").WithCause(ioerror)
//   }
//
func (e *Error) WithCause(cause error) *Error {
	if e.cause == nil {
		e.cause = cause
	}
	return e
}

// supports interface builtin.error
func (e *Error) Error() string {
	var causestr string
	if e.cause != nil {
		causestr = fmt.Sprintf(" (cause: %s)", e.cause.Error())
	}
	return fmt.Sprintf("%s%s", e.error.Error(), causestr)
}

// Returns true if the Error.error is an 'instance'
// of input arg 'errfn'.
func (e *Error) Is(efn errFn) bool {
	s := e.Error()
	category := efn().Error()
	catlen := len(category)
	if len(s) < catlen || s[:catlen] != category {
		return false
	}

	return true
}
