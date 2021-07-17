package stealthpool

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiError(t *testing.T) {
	testCases := []struct {
		name             string
		errs             []error
		shouldHaveErrors bool
		expected         string
	}{
		{
			name:             "no errors",
			shouldHaveErrors: false,
		},
		{
			name:             "nil error",
			errs:             []error{nil},
			shouldHaveErrors: false,
		},
		{
			name:             "nil errors",
			errs:             []error{nil, nil, nil},
			shouldHaveErrors: false,
		},
		{
			name:             "one error",
			errs:             []error{errors.New("Boom")},
			shouldHaveErrors: true,
			expected:         "Boom",
		},
		{
			name:             "multiple errors",
			errs:             []error{errors.New("Boom"), errors.New("shaka"), errors.New("laka")},
			shouldHaveErrors: true,
			expected:         "Boom\nshaka\nlaka",
		},
		{
			name:             "multiple errors with nils",
			errs:             []error{errors.New("Boom"), nil, errors.New("shaka"), nil, errors.New("laka"), nil},
			shouldHaveErrors: true,
			expected:         "Boom\nshaka\nlaka",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			me := newMultiErr()
			for _, err := range tc.errs {
				me.Add(err)
			}
			if tc.shouldHaveErrors {
				require.NotNil(me.Return())
				assert.Equal(tc.expected, me.Error())
			} else {
				require.Nil(me.Return())
			}
		})
	}
}
