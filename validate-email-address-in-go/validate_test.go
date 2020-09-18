package main

import (
	"fmt"
	"testing"
)

// TestEmailValidity is the unit test to test validEmail function.
func TestEmailValidity(t *testing.T) {
	var tests = []struct {
		email string
		want  bool
	}{
		{email: "test@gmail.com", want: true},
		{email: "test@notgmail.com", want: false},
		{email: "test@", want: false},
		{email: "@gmail.com", want: false},
	}

	for _, td := range tests {
		testname := fmt.Sprintf("when testing email id `%s`", td.email)
		t.Run(testname, func(t *testing.T) {
			ok := validEmail(td.email)
			if ok != td.want {
				t.Errorf("got %t, want %t", ok, td.want)
			}
		})
	}
}
