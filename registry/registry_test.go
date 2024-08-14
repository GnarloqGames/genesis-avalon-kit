package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID(t *testing.T) {
	version := "1.0.0"
	resource := ResourceBlueprintRequest{
		Name: "test",
		Slug: "test",
	}

	assert.Equal(t, "d1525252-4413-508b-9c2a-12a96a92353b", ID(resource, version).String())
}
