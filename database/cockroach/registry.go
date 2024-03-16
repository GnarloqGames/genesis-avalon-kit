package cockroach

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"

	"github.com/GnarloqGames/genesis-avalon-kit/proto"
	sq "github.com/Masterminds/squirrel"
)

func (c *Connection) SaveBuildingBlueprint(ctx context.Context, blueprint *proto.BuildingBlueprint) error {
	buf := bytes.NewBufferString("")
	encoder := json.NewEncoder(buf)

	if err := encoder.Encode(blueprint); err != nil {
		return err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, params, err := psql.Insert("building_blueprints").
		Columns("version", "name", "slug", "definition").
		Values(blueprint.Version, blueprint.Name, blueprint.Slug, buf.String()).
		ToSql()

	if err != nil {
		return err
	}

	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, query, params...); err != nil {
		if err := tx.Rollback(ctx); err != nil {
			slog.Error("failed to roll back transaction")
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
