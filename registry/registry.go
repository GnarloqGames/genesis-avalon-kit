package registry

import (
	"time"

	"github.com/GnarloqGames/genesis-avalon-kit/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	Version = "1.0.0"

	SlugWood  = "wood"
	SlugHouse = "house"
)

type Store struct {
	Resources map[string]*proto.ResourceBlueprint
	Buildings map[string]*proto.BuildingBlueprint
}

var (
	wood = &proto.ResourceBlueprint{
		ID:      uuid.NewSHA1(uuid.NameSpaceOID, []byte(SlugWood)).String(),
		Name:    "Wood",
		Slug:    SlugWood,
		Version: Version,
	}

	house = &proto.BuildingBlueprint{
		ID:        uuid.NewSHA1(uuid.NameSpaceOID, []byte(SlugHouse)).String(),
		Name:      "House",
		Slug:      SlugHouse,
		Version:   Version,
		BuildTime: durationpb.New(30 * time.Second),
		Cost: &proto.ResourceList{
			Resources: []*proto.ResourceListItem{
				{
					Resource: wood,
					Amount:   10,
				},
			},
		},
	}
)

var Blueprints = Store{
	Resources: map[string]*proto.ResourceBlueprint{
		SlugWood: wood,
	},

	Buildings: map[string]*proto.BuildingBlueprint{
		SlugHouse: house,
	},
}
