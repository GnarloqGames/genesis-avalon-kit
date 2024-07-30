package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/GnarloqGames/genesis-avalon-kit/database/cockroach/driver"
	"github.com/GnarloqGames/genesis-avalon-kit/proto"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Inventory struct {
	conn *driver.Connection
}

func NewInventory(c *driver.Connection) *Inventory {
	return &Inventory{
		conn: c,
	}
}

func (c *Inventory) GetBuildings(ctx context.Context, owner uuid.UUID) ([]*proto.Building, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, params, err := psql.Select("id", "blueprint", "id_owner", "created_at", "finished_at", "active").
		From("inventory_buildings").
		Where("id_owner = ?", owner).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := c.conn.Query(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer rows.Close()
	results := make([]*proto.Building, 0)
	for rows.Next() {
		result, err := pgx.RowToMap(rows)
		if err != nil {
			return nil, err
		}

		idRaw := result["id"].([16]uint8)
		id, err := uuid.FromBytes(idRaw[:])
		if err != nil {
			return nil, err
		}

		idOwnerRaw := result["id_owner"].([16]uint8)
		idOwner, err := uuid.FromBytes(idOwnerRaw[:])
		if err != nil {
			return nil, err
		}

		building := &proto.Building{
			ID:        id.String(),
			Owner:     idOwner.String(),
			Active:    result["active"].(bool),
			Blueprint: result["blueprint"].(string),
			BuiltAt:   timestamppb.New(result["finished_at"].(time.Time)),
		}

		results = append(results, building)
	}

	return results, nil
}

func (c *Inventory) GetBuilding(ctx context.Context, id uuid.UUID) (*proto.Building, error) {
	return nil, nil
}

func (c *Inventory) SaveBuilding(ctx context.Context, building *proto.Building) error {
	return nil
}
