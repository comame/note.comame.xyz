package test

import (
	"log"
	"reflect"
	"testing"
)

func AssertEquals[T any](t *testing.T, got, expect T) {
	if reflect.DeepEqual(got, expect) {
		return
	}

	log.Printf("got    %v", got)
	log.Printf("expect %v", expect)
	t.Fail()
}

func AssertSame[T comparable](t *testing.T, got, expect T) {
	if got == expect {
		return
	}

	log.Printf("got    %v", got)
	log.Printf("expect %v", expect)
	t.Fail()
}
