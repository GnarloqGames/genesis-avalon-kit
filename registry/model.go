package registry

type Building struct {
	ID      string `json:"id"`
	BuiltAt string `json:"built_at"`
	Active  bool   `json:"active"`
}
