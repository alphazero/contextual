// Copyright 2010-2016 Joubin Houshyar.  All rights reserved.
// Use of this source code is governed by a 2-clause BSD
// license that can be found in the LICENSE file.

// whitebox tests
package goerror_test

import (
	"errors"
	"goerror"
	"testing"
)

var fubarError = goerror.Define("fubar")

func TestDefine(t *testing.T) {
	// both test applicability of TypeOf() for any error
	// and test that match fails
	e0 := errors.New("generic")
	if goerror.TypeOf(e0).Is(fubarError) {
		t.Errorf("e0 is not an fubarError (error)")
	}

	// match the len of the error category string
	// to insure we're not simply matching string lengths
	errPrefix := "error - " // REVU: don't like this in general but in a whitebox test?
	e1 := errors.New(errPrefix + "spoO")
	if goerror.TypeOf(e1).Is(fubarError) {
		t.Errorf("e1 is not an fubarError (error)")
	}

	// use the error as is (no detailed info)
	e := fubarError()
	if !goerror.TypeOf(e).Is(fubarError) {
		t.Errorf("e(no args) should be an fubarError (error)")
	}

	// use the error with a simple string detail msg
	e = fubarError("did it again!")
	if !goerror.TypeOf(e).Is(fubarError) {
		t.Errorf("e should be an fubarError (error)")
	}
}
