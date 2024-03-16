package database

import (
	"context"

	"github.com/GnarloqGames/genesis-avalon-kit/database/cockroach"
	"github.com/GnarloqGames/genesis-avalon-kit/proto"
)

type Store interface {
	// GetBuildings(ctx context.Context, owner uuid.UUID) ([]*proto.Building, error)
	// GetBuilding(ctx context.Context, id uuid.UUID) (*proto.Building, error)
	// SaveBuilding(ctx context.Context, building *proto.Building) error
	SaveBuildingBlueprint(ctx context.Context, blueprint *proto.BuildingBlueprint) error
}

var (
	// _ Store = (*couchbase.Connection)(nil)
	_ Store = (*cockroach.Connection)(nil)
)
