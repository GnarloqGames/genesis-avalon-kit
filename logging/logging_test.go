package logging_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/GnarloqGames/genesis-avalon-kit/logging"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		level         string
		kind          string
		expectedError bool
	}{
		{
			level:         "debug",
			kind:          "json",
			expectedError: false,
		},
		{
			level:         "info",
			kind:          "text",
			expectedError: false,
		},
		{
			level:         "warn",
			kind:          "json",
			expectedError: false,
		},
		{
			level:         "error",
			kind:          "json",
			expectedError: false,
		},
		{
			level:         "debugg",
			kind:          "json",
			expectedError: true,
		},
		{
			level:         "debug",
			kind:          "jsonn",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		tf := func(t *testing.T) {
			_, err := logging.Logger(
				logging.WithKind(tt.kind),
				logging.WithLevel(tt.level),
				logging.WithOutput(os.Stderr),
			)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		}
		t.Run(fmt.Sprintf("%s-%s-%t", tt.level, tt.kind, tt.expectedError), tf)
	}
}
