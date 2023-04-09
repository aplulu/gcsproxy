package accesstoken

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccessToken_Valid(t *testing.T) {
	testCases := []struct {
		name        string
		AccessToken *AccessToken
		wantErr     error
	}{{
		name: "Success",
		AccessToken: &AccessToken{
			Issuer:         "https://example.com",
			ExpirationTime: time.Now().Unix() + 3600,
			Audience:       Audience{"https://example.com"},
			Subject:        "USER_ID",
			IssuedAt:       time.Now().Unix(),
		},
	}, {
		name: "Expired",
		AccessToken: &AccessToken{
			Issuer:         "https://example.com",
			ExpirationTime: time.Now().Unix() - 3600,
			Audience:       Audience{"https://example.com"},
			Subject:        "USER_ID",
			IssuedAt:       time.Now().Unix() - 7200,
		},
		wantErr: ErrExpiredToken,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.AccessToken.Valid()

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
