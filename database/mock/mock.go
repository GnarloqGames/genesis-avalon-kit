package mock

import (
	"context"
	"fmt"

	"github.com/GnarloqGames/genesis-avalon-kit/proto"
	"github.com/google/uuid"
)

var store *Store

type Store struct {
	Buildings []*proto.Building

	ResourceBlueprints []*proto.ResourceBlueprint
	BuildingBlueprints []*proto.BuildingBlueprint

	TaskStatus map[string]*proto.TaskStatus
}

func Get() (*Store, error) {
	if store == nil {
		store = &Store{
			Buildings: make([]*proto.Building, 0),

			ResourceBlueprints: make([]*proto.ResourceBlueprint, 0),
			BuildingBlueprints: make([]*proto.BuildingBlueprint, 0),
		}
	}

	return store, nil
}

func (s *Store) GetBuildings(ctx context.Context, owner uuid.UUID) ([]*proto.Building, error) {
	response := make([]*proto.Building, 0)

	for _, item := range s.Buildings {
		if item.Owner == owner.String() {
			response = append(response, item)
		}
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("not found")
	}

	return response, nil
}

func (s *Store) GetBuilding(ctx context.Context, id uuid.UUID) (*proto.Building, error) {
	for _, item := range s.Buildings {
		if item.ID == id.String() {
			return item, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

func (s *Store) SaveBuilding(ctx context.Context, building *proto.Building) error {
	s.Buildings = append(s.Buildings, building)

	return nil
}

func (s *Store) SaveBuildingBlueprint(ctx context.Context, blueprint *proto.BuildingBlueprint) error {
	s.BuildingBlueprints = append(s.BuildingBlueprints, blueprint)

	return nil
}

func (s *Store) SaveResourceBlueprint(ctx context.Context, blueprint *proto.ResourceBlueprint) error {
	s.ResourceBlueprints = append(s.ResourceBlueprints, blueprint)

	return nil
}

func (s *Store) GetBuildingBlueprints(ctx context.Context, version string) ([]*proto.BuildingBlueprint, error) {
	response := make([]*proto.BuildingBlueprint, 0)

	for _, item := range s.BuildingBlueprints {
		if item.Version == version {
			response = append(response, item)
		}
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("not found")
	}

	return response, nil
}

func (s *Store) GetResourceBlueprints(ctx context.Context, version string) ([]*proto.ResourceBlueprint, error) {
	response := make([]*proto.ResourceBlueprint, 0)

	for _, item := range s.ResourceBlueprints {
		if item.Version == version {
			response = append(response, item)
		}
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("not found")
	}

	return response, nil
}

func (s *Store) GetBuildingBlueprint(ctx context.Context, version string, slug string) (*proto.BuildingBlueprint, error) {
	for _, item := range s.BuildingBlueprints {
		if item.Version == version && item.Slug == slug {
			return item, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

func (s *Store) GetResourceBlueprint(ctx context.Context, version string, slug string) (*proto.ResourceBlueprint, error) {
	for _, item := range s.ResourceBlueprints {
		if item.Version == version && item.Slug == slug {
			return item, nil
		}
	}

	return nil, fmt.Errorf("not found")
}

func (s *Store) SaveTaskStatus(ctx context.Context, status *proto.TaskStatus) error {
	s.TaskStatus[status.ID] = status
	return nil
}
