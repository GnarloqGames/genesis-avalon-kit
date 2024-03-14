package cockroach

import "fmt"

var (
	ErrEmptyHostname = fmt.Errorf("hostname cannot be empty")
	ErrEmptyPort     = fmt.Errorf("port cannot be 0")
	ErrEmptyDatabase = fmt.Errorf("database cannot be empty")
)
