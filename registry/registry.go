package registry

import (
	"context"
	"fmt"
	"time"

	"github.com/GnarloqGames/genesis-avalon-kit/database"
	"github.com/GnarloqGames/genesis-avalon-kit/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	Version = "1.0.0"

	SlugWood  = "wood"
	SlugHouse = "house"
)

var (
	wood = &proto.ResourceBlueprint{
		ID:      uuid.NewSHA1(uuid.NameSpaceOID, []byte(fmt.Sprintf("%s:%s", SlugWood, Version))).String(),
		Name:    "Wood",
		Slug:    SlugWood,
		Version: Version,
	}

	house = &proto.BuildingBlueprint{
		ID:        uuid.NewSHA1(uuid.NameSpaceOID, []byte(fmt.Sprintf("%s:%s", SlugHouse, Version))).String(),
		Name:      "House",
		Slug:      SlugHouse,
		Version:   Version,
		BuildTime: durationpb.New(30 * time.Second),
		Cost: &proto.ResourceList{
			Resources: []*proto.ResourceListItem{
				{
					Name:   SlugWood,
					Amount: 10,
				},
			},
		},
	}
)

var Blueprints = Store{
	Resources: NewItemStore[*proto.ResourceBlueprint](),
	Buildings: NewItemStore[*proto.BuildingBlueprint](),
}

func init() {
	Blueprints.Resources.Set(SlugWood, wood)
	Blueprints.Buildings.Set(SlugHouse, house)
}

func SaveBuildingBlueprint(ctx context.Context, blueprint BuildingBlueprintRequest) error {
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
		ID:         ID(blueprint).String(),
		Name:       blueprint.Name,
		Slug:       blueprint.Slug,
		Version:    blueprint.Version,
		BuildTime:  durationpb.New(duration),
		Cost:       &proto.ResourceList{Resources: resources},
		Production: production,
	}

	return conn.SaveBuildingBlueprint(ctx, blueprintProto)
}

func SaveResourceBlueprint(ctx context.Context, resource ResourceBlueprintRequest) error {
	conn, err := database.Get()
	if err != nil {
		return err
	}

	resourceProto := &proto.ResourceBlueprint{
		ID:      ID(resource).String(),
		Name:    resource.Name,
		Slug:    resource.Slug,
		Version: resource.Version,
	}

	return conn.SaveResourceBlueprint(ctx, resourceProto)
}

func LoadBuildingBlueprints(ctx context.Context, version string) error {
	conn, err := database.Get()
	if err != nil {
		return err
	}

	buildings, err := conn.GetBuildingBlueprints(ctx, version)
	if err != nil {
		return err
	}

	for _, building := range buildings {
		store.Buildings.Set(building.Slug, building)
	}

	return nil
}

func LoadResourceBlueprints(ctx context.Context, version string) error {
	conn, err := database.Get()
	if err != nil {
		return err
	}

	resources, err := conn.GetResourceBlueprints(ctx, version)
	if err != nil {
		return err
	}

	for _, resource := range resources {
		store.Resources.Set(resource.Slug, resource)
	}

	return nil
}
