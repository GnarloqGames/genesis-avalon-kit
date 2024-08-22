package cache

import (
	"context"
	"fmt"
	"testing"

	"github.com/GnarloqGames/genesis-avalon-kit/database"
	"github.com/GnarloqGames/genesis-avalon-kit/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errNotFound = fmt.Errorf("not found")
)

func TestLoad(t *testing.T) {
	database.SetKind(database.DriverMock)

	expectedBuildingBlueprint := &proto.BuildingBlueprint{
		Name:    "Test",
		Slug:    "test",
		Version: "1.0.0",
	}

	expectedResourceBlueprint := &proto.ResourceBlueprint{
		Name:    "Test",
		Slug:    "test",
		Version: "1.0.0",
	}

	mock, _ := database.Get() //nolint
	mock.SaveBuildingBlueprint(context.Background(), expectedBuildingBlueprint)
	mock.SaveResourceBlueprint(context.Background(), expectedResourceBlueprint)

	tests := []struct {
		label         string
		version       string
		expectedError error
	}{
		{
			label:         "success",
			version:       "1.0.0",
			expectedError: nil,
		},
		{
			label:         "failed-resource",
			version:       "0.0.1",
			expectedError: errNotFound,
		},
		{
			label:         "failed-building",
			version:       "0.0.2",
			expectedError: errNotFound,
		},
	}

	for _, tt := range tests {
		tf := func(t *testing.T) {
			store = nil
			SetVersion(tt.version)
			err := Load(context.Background())

			if tt.expectedError != nil {
				assert.ErrorContains(t, err, "not found")
			} else {
				require.NoError(t, err)

				resource, ok := store.Resources.Get("test")
				require.True(t, ok)
				assert.Equal(t, "Test", resource.Name)
				assert.Equal(t, "test", resource.Slug)
				assert.Equal(t, "1.0.0", resource.Version)

				building, ok := store.Buildings.Get("test")
				require.True(t, ok)
				assert.Equal(t, "Test", building.Name)
				assert.Equal(t, "test", building.Slug)
				assert.Equal(t, "1.0.0", building.Version)
			}
		}

		t.Run(tt.version, tf)
	}
}
