package cockroach

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

var conn *Connection

type Connection struct {
	*pgx.Conn
}

func Get() (*Connection, error) {
	if conn == nil {
		dsn, err := connectionString()
		if err != nil {
			return nil, err
		}

		ctx := context.Background()
		c, err := pgx.Connect(ctx, dsn)
		if err != nil {
			return nil, err
		}

		conn = &Connection{c}
	}

	return conn, nil
}

func connectionString() (string, error) {
	if hostname == "" {
		return "", ErrEmptyHostname
	}

	if port == 0 {
		return "", ErrEmptyPort
	}

	if database == "" {
		return "", ErrEmptyDatabase
	}

	dsn := fmt.Sprintf("%s:%d/%s?sslmode=verify-full", hostname, port, database)

	if username != "" {
		auth := username
		if password != "" {
			auth = fmt.Sprintf("%s:%s", auth, password)
		}
		dsn = fmt.Sprintf("%s@%s", auth, dsn)
	}

	return fmt.Sprintf("postgresql://%s", dsn), nil
}
