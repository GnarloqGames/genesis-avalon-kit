package couchbase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/GnarloqGames/genesis-avalon-kit/proto"
	"github.com/GnarloqGames/genesis-avalon-kit/registry"
	"github.com/couchbase/gocb/v2"
	"github.com/couchbase/gocbcore/v10"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			RetryStrategy: gocb.NewBestEffortRetryStrategy(backoffCalculator()),
		})

		if err != nil {
			return nil, err
		}

		b := cluster.Bucket(Bucket())
		if err := b.WaitUntilReady(15*time.Second, nil); err != nil {
			return nil, err
		}

		conn = &Connection{cluster}
	}

	return conn, nil
}

func (c *Connection) upsert(item any) error {
	conn, err := Get()
	if err != nil {
		return err
	}

	bucket := conn.Bucket(Bucket())

	if err := bucket.WaitUntilReady(3*time.Second, nil); err != nil {
		return err
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

func (c *Connection) GetBuildings(ctx context.Context, owner uuid.UUID) ([]*proto.Building, error) {
	bucket := conn.Bucket(Bucket())

	if err := bucket.WaitUntilReady(3*time.Second, nil); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT * FROM buildings WHERE owner = "%s"`, owner.String())
	queryResult, err := bucket.Scope(Scope()).Query(query, nil)
	if err != nil {
		return nil, err
	}

	results := make([]*proto.Building, 0)
	for queryResult.Next() {
		var result map[string]interface{}
		if err := queryResult.Row(&result); err != nil {
			return nil, err
		}

		entity := result["buildings"].(map[string]interface{})

		buildingName := entity["name"].(string)
		blueprint, ok := registry.Blueprints.Buildings[buildingName]
		if !ok {
			return nil, fmt.Errorf("invalid building type %s", buildingName)
		}

		builtAt, err := time.Parse(time.RFC1123, entity["built_at"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to parse building timestamp: %w", err)
		}

		building := &proto.Building{
			ID:        entity["id"].(string),
			Blueprint: blueprint,
			BuiltAt:   timestamppb.New(builtAt),
			Active:    entity["active"].(bool),
		}

		results = append(results, building)
	}

	return results, nil
}

func (c *Connection) GetBuilding(ctx context.Context, id uuid.UUID) (*proto.Building, error) {
	bucket := conn.Bucket(Bucket())

	if err := bucket.WaitUntilReady(3*time.Second, nil); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`SELECT * FROM buildings WHERE id = "%s"`, id.String())
	queryResult, err := bucket.Scope(Scope()).Query(query, nil)
	if err != nil {
		return nil, err
	}

	queryResult.Next()
	var result map[string]interface{}
	if err := queryResult.Row(&result); err != nil {
		return nil, err
	}

	entity := result["buildings"].(map[string]interface{})

	buildingName := entity["name"].(string)
	blueprint, ok := registry.Blueprints.Buildings[buildingName]
	if !ok {
		return nil, fmt.Errorf("invalid building type %s", buildingName)
	}

	builtAt, err := time.Parse(time.RFC1123, entity["built_at"].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to parse building timestamp: %w", err)
	}

	building := &proto.Building{
		ID:        entity["id"].(string),
		Blueprint: blueprint,
		BuiltAt:   timestamppb.New(builtAt),
		Active:    entity["active"].(bool),
	}

	return building, nil
}

func (c *Connection) SaveBuilding(ctx context.Context, building *proto.Building) error {
	return c.upsert(building)
}

func backoffCalculator() gocb.BackoffCalculator {
	return gocb.BackoffCalculator(gocbcore.ExponentialBackoff(1*time.Millisecond, 5*time.Second, 2))
}
