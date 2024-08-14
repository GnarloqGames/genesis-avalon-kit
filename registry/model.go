package registry

import (
	"fmt"

	"github.com/google/uuid"
)

type Building struct {
	ID      string `json:"id"`
	BuiltAt string `json:"built_at"`
	Active  bool   `json:"active"`
}

func Hash(name string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(name))
}

type BuildingBlueprintRequest struct {
	Name       string       `json:"name"`
	Slug       string       `json:"slug"`
	BuildTime  string       `json:"build_time"`
	Cost       ResourceList `json:"cost"`
	Production []Production `json:"production"`
}

func (r BuildingBlueprintRequest) GetName() string { return r.Name }
func (r BuildingBlueprintRequest) GetSlug() string { return r.Slug }

type ResourceList []ResourceListItem

type ResourceListItem struct {
	Resource string `json:"resource"`
	Amount   uint64 `json:"amount"`
}

type Production struct {
	Cost           ResourceList `json:"cost"`
	Product        ResourceList `json:"product"`
	ProductionTime string       `json:"production_time"`
}

type ResourceBlueprintRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (r ResourceBlueprintRequest) GetName() string { return r.Name }
func (r ResourceBlueprintRequest) GetSlug() string { return r.Slug }

type Request interface {
	GetName() string
	GetSlug() string
}

func ID(r Request, version string) uuid.UUID {
	id := fmt.Sprintf("%s:%s", r.GetSlug(), version)
	hash := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(id))

	return hash
}
