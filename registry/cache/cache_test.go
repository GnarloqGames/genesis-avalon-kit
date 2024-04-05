package cache

import (
	"context"
	"fmt"
	"testing"

	"github.com/GnarloqGames/genesis-avalon-kit/proto"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errFailedBuilding = fmt.Errorf("failed building")
	errFailedResource = fmt.Errorf("failed resource")
)

func TestLoad(t *testing.T) {
	patches := gomonkey.ApplyFunc(LoadResourceBlueprints, func(ctx context.Context, version string) error {
		if version == "0.0.1" {
			return errFailedResource
		}

		store.Resources.Set("test", &proto.ResourceBlueprint{
			ID:      "1",
			Name:    "Test",
			Slug:    "test",
			Version: "1.0.0",
		})
		return nil
	})
	patches.ApplyFunc(LoadBuildingBlueprints, func(ctx context.Context, version string) error {
		if version == "0.0.2" {
			return errFailedBuilding
		}

		store.Buildings.Set("test", &proto.BuildingBlueprint{
			ID:      "1",
			Name:    "Test",
			Slug:    "test",
			Version: "1.0.0",
		})
		return nil
	})

	defer patches.Reset()

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
			expectedError: errFailedResource,
		},
		{
			label:         "failed-building",
			version:       "0.0.2",
			expectedError: errFailedBuilding,
		},
	}

	for _, tt := range tests {
		tf := func(t *testing.T) {
			store = nil
			SetVersion(tt.version)
			err := Load(context.Background())

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
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
