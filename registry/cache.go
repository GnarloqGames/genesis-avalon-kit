package registry

import (
	"context"
	"sync"

	"github.com/GnarloqGames/genesis-avalon-kit/proto"
)

var store *Store

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
	store = &Store{
		Resources: NewItemStore[*proto.ResourceBlueprint](),
		Buildings: NewItemStore[*proto.BuildingBlueprint](),
	}

	if err := LoadResourceBlueprints(ctx, version); err != nil {
		return err
	}

	if err := LoadBuildingBlueprints(ctx, version); err != nil {
		return err
	}

	return nil
}

func GetLoadedBlueprints(ctx context.Context) map[string]any {
	s := make(map[string]any)

	s["resources"] = store.Resources.GetItems()
	s["buildings"] = store.Buildings.GetItems()

	return s
}
