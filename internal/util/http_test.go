package util

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTLS(t *testing.T) {
	testCases := []struct {
		name string
		arg  *http.Request
		want bool
	}{{
		name: "Is TLS #1",
		arg: &http.Request{
			TLS: &tls.ConnectionState{},
		},
		want: true,
	}, {
		name: "Is TLS #2",
		arg: &http.Request{
			Header: map[string][]string{
				"X-Forwarded-Proto": {"https"},
			},
		},
		want: true,
	}, {
		name: "Is not TLS #1",
		arg:  &http.Request{},
	}, {
		name: "Is not TLS #2",
		arg: &http.Request{
			Header: map[string][]string{
				"X-Forwarded-Proto": {"http"},
			},
		},
	}, {
		name: "Is not TLS #3",
		arg: &http.Request{
			Header: map[string][]string{
				"X-Forwarded-Proto": {},
			},
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsTLS(tc.arg)

			assert.Equal(t, tc.want, got)
		})
	}
}
