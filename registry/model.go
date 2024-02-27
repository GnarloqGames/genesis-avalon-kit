package registry

type Building struct {
	ID     string `json:"id"`
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
