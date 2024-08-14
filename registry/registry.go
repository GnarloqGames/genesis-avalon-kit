package registry

import (
	"context"
	"time"

	"github.com/GnarloqGames/genesis-avalon-kit/database"
	"github.com/GnarloqGames/genesis-avalon-kit/proto"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	Version = "1.0.0"

	SlugWood  = "wood"
	SlugHouse = "house"
)

func SaveBuildingBlueprint(ctx context.Context, version string, blueprint BuildingBlueprintRequest, force bool) error {
	conn, err := database.Get()
	if err != nil {
		return err
	}

	resources := make([]*proto.ResourceListItem, 0)
	for _, item := range blueprint.Cost {
		resources = append(resources, &proto.ResourceListItem{
			Name:   item.Resource,
			Amount: item.Amount,
		})
	}

	production := make([]*proto.Production, 0)
	for _, item := range blueprint.Production {
		cost := make([]*proto.ResourceListItem, 0)
		for _, costItem := range item.Cost {
			cost = append(cost, &proto.ResourceListItem{
				Name:   costItem.Resource,
				Amount: costItem.Amount,
			})
		}

		product := make([]*proto.ResourceListItem, 0)
		for _, productItem := range item.Product {
			product = append(product, &proto.ResourceListItem{
				Name:   productItem.Resource,
				Amount: productItem.Amount,
			})
		}

		productionTime, err := time.ParseDuration(item.ProductionTime)
		if err != nil {
			return err
		}

		production = append(production, &proto.Production{
			Cost:           &proto.ResourceList{Resources: cost},
			Output:         &proto.ResourceList{Resources: product},
			ProductionTime: durationpb.New(productionTime),
		})
	}

	duration, err := time.ParseDuration(blueprint.BuildTime)
	if err != nil {
		return err
	}

	blueprintProto := &proto.BuildingBlueprint{
		ID:         ID(blueprint, version).String(),
		Name:       blueprint.Name,
		Slug:       blueprint.Slug,
		Version:    version,
		BuildTime:  durationpb.New(duration),
		Cost:       &proto.ResourceList{Resources: resources},
		Production: production,
	}

	return conn.SaveBuildingBlueprint(ctx, blueprintProto)
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

func GetBuildingBlueprint(ctx context.Context, version string, slug string) (*proto.BuildingBlueprint, error) {
	store, err := database.Get()
	if err != nil {
		return nil, err
	}

	pb, err := store.GetBuildingBlueprint(ctx, version, slug)
	if err != nil {
		return nil, err
	}

	return pb, nil
}

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
