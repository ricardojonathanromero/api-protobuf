package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		title  string
		setEnv bool
		name   string
		value  string
	}{
		{
			title:  "happy_path",
			setEnv: true,
			name:   "EXAMPLE",
			value:  "some-value",
		},
		{
			title:  "empty_response",
			setEnv: true,
			name:   "EMPTY",
			value:  "",
		},
		{
			title:  "default_value",
			setEnv: false,
			name:   "NOT_EXISTS",
			value:  "default",
		},
	}

	for _, test := range tests {
		if test.setEnv {
			_ = os.Setenv(test.name, test.value)
		}

		t.Run(test.title, func(t *testing.T) {
			returned := GetEnv(test.name, test.value)

			t.Logf("returned [%s] - expected [%s]", returned, test.value)
			if !assert.Equal(t, test.value, returned) {
				t.Errorf("environment %s is not equal to %s", test.name, test.value)
				t.FailNow()
			}
		})

		if test.setEnv {
			_ = os.Remove(test.name)
		}
	}
}

func TestContextWithTimeout(t *testing.T) {
	tests := []struct {
		title   string
		seconds int
	}{
		{
			title:   "happy_path",
			seconds: 10,
		},
		{
			title:   "apply_default_timeout",
			seconds: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			_, cancel := ContextWithTimeout(test.seconds)

			cancel()
		})
	}
}
