package registry

import (
	"context"

	"github.com/GnarloqGames/genesis-avalon-kit/database"
	"github.com/GnarloqGames/genesis-avalon-kit/proto"
)

func GetResourceBlueprint(ctx context.Context, version string, slug string) (*proto.ResourceBlueprint, error) {
	store, err := database.Get()
	if err != nil {
		return nil, err
	}

	pb, err := store.GetResourceBlueprint(ctx, version, slug)
	if err != nil {
		return nil, err
	}

	return pb, nil
}

func SaveResourceBlueprint(ctx context.Context, version string, resource ResourceBlueprintRequest, force bool) error {
	conn, err := database.Get()
	if err != nil {
		return err
	}

	resourceProto := &proto.ResourceBlueprint{
		ID:      ID(resource, version).String(),
		Name:    resource.Name,
		Slug:    resource.Slug,
		Version: version,
	}

	return conn.SaveResourceBlueprint(ctx, resourceProto)
}
