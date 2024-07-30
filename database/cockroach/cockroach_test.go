package cockroach

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectionString(t *testing.T) {
	tests := []struct {
		label       string
		username    string
		password    string
		database    string
		port        uint16
		hostname    string
		expectedStr string
		expectedErr error
	}{
		{
			label:       "happy-path",
			username:    "foo",
			password:    "bar",
			database:    "test",
			port:        26257,
			hostname:    "localhost",
			expectedStr: "postgresql://foo:bar@localhost:26257/test?sslmode=verify-full",
		},
		{
			label:       "only-user",
			username:    "foo",
			password:    "",
			database:    "test",
			port:        26257,
			hostname:    "localhost",
			expectedStr: "postgresql://foo@localhost:26257/test?sslmode=verify-full",
		},
		{
			label:       "no-auth",
			username:    "",
			password:    "",
			database:    "test",
			port:        26257,
			hostname:    "localhost",
			expectedStr: "postgresql://localhost:26257/test?sslmode=verify-full",
		},
		{
			label:       "missing-db",
			username:    "foo",
			password:    "bar",
			database:    "",
			port:        26257,
			hostname:    "localhost",
			expectedErr: ErrEmptyDatabase,
		},
		{
			label:       "missing-hostname",
			username:    "foo",
			password:    "bar",
			database:    "test",
			port:        26257,
			hostname:    "",
			expectedErr: ErrEmptyHostname,
		},
		{
			label:       "missing-port",
			username:    "foo",
			password:    "bar",
			database:    "test",
			port:        0,
			hostname:    "localhost",
			expectedErr: ErrEmptyPort,
		},
	}

	for _, tt := range tests {
		tf := func(t *testing.T) {
			SetHostname(tt.hostname)
			SetUsername(tt.username)
			SetPassword(tt.password)
			SetPort(tt.port)
			SetDatabase(tt.database)

			str, err := connectionString()

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			}

			if tt.expectedStr != "" {
				assert.Equal(t, tt.expectedStr, str)
			}
		}

		t.Run(tt.label, tf)
	}
}
