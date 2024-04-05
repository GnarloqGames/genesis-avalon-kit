package cockroach

import (
	"context"
	"fmt"

	"github.com/GnarloqGames/genesis-avalon-kit/database/cockroach/driver"
	"github.com/GnarloqGames/genesis-avalon-kit/database/cockroach/repository"
	"github.com/jackc/pgx/v5"
)

var backend *Backend

type Backend struct {
	*repository.Registry
	*repository.Inventory
}

func Get() (*Backend, error) {
	if backend == nil {
		dsn, err := connectionString()
		if err != nil {
			return nil, err
		}

		ctx := context.Background()
		c, err := pgx.Connect(ctx, dsn)
		if err != nil {
			return nil, err
		}

		conn := driver.Connection{Conn: c}

		backend = &Backend{
			Inventory: repository.NewInventory(&conn),
			Registry:  repository.NewRegistry(&conn),
		}
	}

	return backend, nil
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
