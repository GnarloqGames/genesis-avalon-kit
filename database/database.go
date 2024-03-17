package database

import (
	"context"
	"fmt"

	"github.com/GnarloqGames/genesis-avalon-kit/database/cockroach"
	"github.com/GnarloqGames/genesis-avalon-kit/proto"
)

type Store interface {
	// GetBuildings(ctx context.Context, owner uuid.UUID) ([]*proto.Building, error)
	// GetBuilding(ctx context.Context, id uuid.UUID) (*proto.Building, error)
	// SaveBuilding(ctx context.Context, building *proto.Building) error
	SaveBuildingBlueprint(ctx context.Context, blueprint *proto.BuildingBlueprint) error
	SaveResourceBlueprint(ctx context.Context, blueprint *proto.ResourceBlueprint) error

	GetBuildingBlueprints(ctx context.Context, version string) ([]*proto.BuildingBlueprint, error)
	GetResourceBlueprints(ctx context.Context, version string) ([]*proto.ResourceBlueprint, error)
}

var (
	// _ Store = (*couchbase.Connection)(nil)
	_ Store = (*cockroach.Connection)(nil)
)

func Get() (Store, error) {
	switch kind {
	case DriverCockroach:
		return cockroach.Get()
	default:
		return nil, fmt.Errorf("invalid database driver")
	}
}
