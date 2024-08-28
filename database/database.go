package database

import (
	"context"
	"fmt"

	"github.com/GnarloqGames/genesis-avalon-kit/database/cockroach"
	"github.com/GnarloqGames/genesis-avalon-kit/database/mock"
	"github.com/GnarloqGames/genesis-avalon-kit/proto"
	"github.com/google/uuid"
)

type Store interface {
	GetBuildings(ctx context.Context, owner uuid.UUID) ([]*proto.Building, error)
	GetBuilding(ctx context.Context, id uuid.UUID) (*proto.Building, error)

	SaveBuilding(ctx context.Context, building *proto.Building) error

	SaveBuildingBlueprint(ctx context.Context, blueprint *proto.BuildingBlueprint) error
	SaveResourceBlueprint(ctx context.Context, blueprint *proto.ResourceBlueprint) error
	GetBuildingBlueprints(ctx context.Context, version string) ([]*proto.BuildingBlueprint, error)
	GetResourceBlueprints(ctx context.Context, version string) ([]*proto.ResourceBlueprint, error)
	GetBuildingBlueprint(ctx context.Context, version string, slug string) (*proto.BuildingBlueprint, error)
	GetResourceBlueprint(ctx context.Context, version string, slug string) (*proto.ResourceBlueprint, error)

	SaveTaskStatus(ctx context.Context, status *proto.TaskStatus) error
}

var (
	// _ Store = (*couchbase.Connection)(nil)
	_ Store = (*cockroach.Backend)(nil)
)

func Get() (Store, error) {
	switch kind {
	case DriverCockroach:
		return cockroach.Get()
	case DriverMock:
		return mock.Get()
	default:
		return nil, fmt.Errorf("invalid database driver")
	}
}
