package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/GnarloqGames/genesis-avalon-kit/proto"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (c *Registry) GetBuildingBlueprint(ctx context.Context, version string, slug string) (*proto.BuildingBlueprint, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, params, err := psql.Select("id", "version", "name", "slug", "definition").
		From("building_blueprints").
		Where("version = ? AND slug = ?", version, slug).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := c.conn.QueryRow(ctx, query, params...)
	var bp *proto.BuildingBlueprint
	if err := row.Scan(&bp); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return bp, nil
}

func (c *Registry) GetBuildingBlueprints(ctx context.Context, version string) ([]*proto.BuildingBlueprint, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, params, err := psql.Select("id", "version", "name", "slug", "definition").
		From("building_blueprints").
		Where("version = ?", version).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := c.conn.Query(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer rows.Close()
	results := make([]*proto.BuildingBlueprint, 0)
	for rows.Next() {
		blueprintMap, err := pgx.RowToMap(rows)
		if err != nil {
			return nil, err
		}

		var blueprint proto.BuildingBlueprint

		raw, err := json.Marshal(blueprintMap["definition"])
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(raw, &blueprint); err != nil {
			return nil, err
		}

		results = append(results, &blueprint)
	}

	return results, nil
}

func (c *Registry) SaveBuildingBlueprint(ctx context.Context, blueprint *proto.BuildingBlueprint) error {
	buf := bytes.NewBufferString("")
	encoder := json.NewEncoder(buf)

	if err := encoder.Encode(blueprint); err != nil {
		return err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, params, err := psql.Insert("building_blueprints").
		Columns("id", "version", "name", "slug", "definition").
		Values(blueprint.ID, blueprint.Version, blueprint.Name, blueprint.Slug, buf.String()).
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
