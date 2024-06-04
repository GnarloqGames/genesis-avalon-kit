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
	Version    string       `json:"version"`
	BuildTime  string       `json:"build_time"`
	Cost       ResourceList `json:"cost"`
	Production []Production `json:"production"`
}

func (r BuildingBlueprintRequest) GetName() string    { return r.Name }
func (r BuildingBlueprintRequest) GetSlug() string    { return r.Slug }
func (r BuildingBlueprintRequest) GetVersion() string { return r.Version }

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
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Version string `json:"version"`
}

func (r ResourceBlueprintRequest) GetName() string    { return r.Name }
func (r ResourceBlueprintRequest) GetSlug() string    { return r.Slug }
func (r ResourceBlueprintRequest) GetVersion() string { return r.Version }

type Request interface {
	GetName() string
	GetSlug() string
	GetVersion() string
}

func ID(r Request) uuid.UUID {
	id := fmt.Sprintf("%s:%s", r.GetSlug(), r.GetVersion())
	hash := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(id))

	return hash
}
