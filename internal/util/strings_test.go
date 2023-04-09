package util

import "testing"

func TestRandomString(t *testing.T) {
	testCases := []struct {
		name string
		arg  int
	}{{
		name: "length 10",
		arg:  10,
	}, {
		name: "length 20",
		arg:  20,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := RandomString(tc.arg)
			if len(s) != tc.arg {
				t.Errorf("length of string is not %d", tc.arg)
			}
		})
	}
}
