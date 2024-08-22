package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/GnarloqGames/genesis-avalon-kit/database"
	"github.com/GnarloqGames/genesis-avalon-kit/proto"
)

var store *Store

type Kind string

const (
	KindBuilding Kind = "building"
	KindResource Kind = "resource"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

type Store struct {
	Resources *ItemStore[*proto.ResourceBlueprint]
	Buildings *ItemStore[*proto.BuildingBlueprint]
}

type Blueprint interface {
	*proto.ResourceBlueprint | *proto.BuildingBlueprint
}

func NewItemStore[T Blueprint]() *ItemStore[T] {
	s := ItemStore[T]{
		mx: &sync.Mutex{},

		items: make(map[string]T),
	}

	return &s
}

type ItemStore[T Blueprint] struct {
	mx *sync.Mutex `json:"-"`

	items map[string]T
}

func (s *ItemStore[T]) GetItems() map[string]T {
	s.mx.Lock()
	defer s.mx.Unlock()

	return s.items
}

func (s *ItemStore[T]) Set(key string, item T) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.items[key] = item
}

func (s *ItemStore[T]) Get(key string) (T, bool) {
	s.mx.Lock()
	defer s.mx.Unlock()

	item, ok := s.items[key]
	return item, ok
}

func Load(ctx context.Context) error {
	newStore := &Store{
		Resources: NewItemStore[*proto.ResourceBlueprint](),
		Buildings: NewItemStore[*proto.BuildingBlueprint](),
	}

	// Load resource blueprints
	if err := LoadBlueprints(ctx, KindResource, newStore, version); err != nil {
		return err
	}

	// Load building blueprints
	if err := LoadBlueprints(ctx, KindBuilding, newStore, version); err != nil {
		return err
	}

	store = newStore

	return nil
}

func GetLoadedBlueprints(ctx context.Context) map[string]any {
	s := make(map[string]any)

	s["resources"] = store.Resources.GetItems()
	s["buildings"] = store.Buildings.GetItems()

	return s
}

func LoadBlueprints(ctx context.Context, kind Kind, store *Store, version string) error {
	conn, err := database.Get()
	if err != nil {
		return err
	}

	switch kind {
	case KindBuilding:
		buildings, err := conn.GetBuildingBlueprints(ctx, version)
		if err != nil {
			return err
		}

		for _, building := range buildings {
			store.Buildings.Set(building.Slug, building)
		}
	case KindResource:
		resources, err := conn.GetResourceBlueprints(ctx, version)
		if err != nil {
			return err
		}

		for _, resource := range resources {
			store.Resources.Set(resource.Slug, resource)
		}
	}

	return nil
}

func GetBuildingBlueprint(ctx context.Context, slug string) (*proto.BuildingBlueprint, bool) {
	return store.Buildings.Get(slug)
}

func GetResourceBlueprint(ctx context.Context, slug string) (*proto.ResourceBlueprint, bool) {
	return store.Resources.Get(slug)
}
