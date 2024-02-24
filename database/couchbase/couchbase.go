package couchbase

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/GnarloqGames/genesis-avalon-kit/registry"
	"github.com/couchbase/gocb/v2"
)

var conn *Connection

type Connection struct {
	*gocb.Cluster
}

func Get() (*Connection, error) {
	if conn == nil {
		connUrl := fmt.Sprintf("couchbase://%s", URL())

		slog.Info("connecting to couchbase", "url", connUrl, "username", username, "bucket", bucket)
		cluster, err := gocb.Connect(connUrl, gocb.ClusterOptions{
			TimeoutsConfig: gocb.TimeoutsConfig{
				ConnectTimeout: 30 * time.Second,
			},
			Authenticator: gocb.PasswordAuthenticator{
				Username: Username(),
				Password: Password(),
			},
		})

		if err != nil {
			return nil, err
		}

		b := cluster.Bucket(Bucket())
		if err := b.WaitUntilReady(3*time.Second, nil); err != nil {
			return nil, err
		}

		conn = &Connection{cluster}
	}

	return conn, nil
}

func (c *Connection) Upsert(item any) error {
	conn, err := Get()
	if err != nil {
		return err
	}

	bucket := conn.Bucket(Bucket())

	if err := bucket.WaitUntilReady(3*time.Second, nil); err != nil {
		log.Fatalln(err)
	}

	collName := ""
	id := ""
	kind := ""

	switch i := item.(type) {
	case registry.Building:
		collName = "buildings"
		id = fmt.Sprintf("b:%s", i.ID)
		kind = "building"
	default:
		return ErrInvalidItemType
	}

	scope := Scope()
	coll := bucket.Scope(scope).Collection(collName)

	slog.Info("upserting item", "item", item, "scope", scope, "collection", collName)

	if _, err := coll.Upsert(id, item, &gocb.UpsertOptions{Timeout: 5 * time.Second}); err != nil {
		// should return error type and log outside
		slog.Error("failed to upsert item", "kind", kind, "item", item)
		return fmt.Errorf("failed to upsert item: %w", err)
	}

	return nil
}
