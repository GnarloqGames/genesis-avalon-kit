package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/GnarloqGames/genesis-avalon-kit/proto"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (c *Registry) GetResourceBlueprint(ctx context.Context, version string, slug string) (*proto.ResourceBlueprint, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, params, err := psql.Select("id", "version", "name", "slug").
		From("resource_blueprints").
		Where("version = ? AND slug = ?", version, slug).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := c.conn.QueryRow(ctx, query, params...)
	var bp *proto.ResourceBlueprint
	if err := row.Scan(&bp); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}

	return bp, nil
}

func (c *Registry) GetResourceBlueprints(ctx context.Context, version string) ([]*proto.ResourceBlueprint, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, params, err := psql.Select("id", "version", "name", "slug").
		From("resource_blueprints").
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
	results := make([]*proto.ResourceBlueprint, 0)
	for rows.Next() {
		blueprint, err := pgx.RowToStructByName[proto.ResourceBlueprint](rows)
		if err != nil {
			return nil, err
		}

		results = append(results, &blueprint)
	}

	return results, nil
}

func (c *Registry) SaveResourceBlueprint(ctx context.Context, blueprint *proto.ResourceBlueprint) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, params, err := psql.Insert("resource_blueprints").
		Columns("id", "version", "name", "slug").
		Values(blueprint.ID, blueprint.Version, blueprint.Name, blueprint.Slug).
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
