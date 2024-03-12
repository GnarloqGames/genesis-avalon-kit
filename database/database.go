package database

import (
	"context"

	"github.com/GnarloqGames/genesis-avalon-kit/database/couchbase"
	"github.com/GnarloqGames/genesis-avalon-kit/proto"
	"github.com/google/uuid"
)

type Store interface {
	GetBuildings(ctx context.Context, owner uuid.UUID) ([]*proto.Building, error)
	GetBuilding(ctx context.Context, id uuid.UUID) (*proto.Building, error)
	SaveBuilding(ctx context.Context, building *proto.Building) error
}

var (
	_ Store = (*couchbase.Connection)(nil)
)
