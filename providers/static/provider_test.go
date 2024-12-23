package static

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	vshardrouter "github.com/tarantool/go-vshard-router"
)

func TestNewProvider(t *testing.T) {
	tCases := []struct {
		Source map[vshardrouter.ReplicasetInfo][]vshardrouter.InstanceInfo
	}{
		{nil},
		{make(map[vshardrouter.ReplicasetInfo][]vshardrouter.InstanceInfo)},
		{map[vshardrouter.ReplicasetInfo][]vshardrouter.InstanceInfo{
			{}: {
				vshardrouter.InstanceInfo{},
				vshardrouter.InstanceInfo{},
			},
		}},
	}

	for _, tc := range tCases {
		tc := tc

		t.Run("provider", func(t *testing.T) {
			t.Parallel()
			if len(tc.Source) == 0 {

				require.Panics(t, func() {
					NewProvider(tc.Source)
				})

				return
			}

			require.NotPanics(t, func() {
				provider := NewProvider(tc.Source)
				require.NotNil(t, provider)
			})
		})
	}
}

func TestProvider_Validate(t *testing.T) {
	tCases := []struct {
		Name   string
		Source map[vshardrouter.ReplicasetInfo][]vshardrouter.InstanceInfo
		IsErr  bool
	}{
		{
			Name: "empty name",
			Source: map[vshardrouter.ReplicasetInfo][]vshardrouter.InstanceInfo{
				{}: {
					vshardrouter.InstanceInfo{},
					vshardrouter.InstanceInfo{},
				},
			},
			IsErr: true,
		},
		{
			Name: "no uuid",
			Source: map[vshardrouter.ReplicasetInfo][]vshardrouter.InstanceInfo{
				{Name: "rs_1"}: {
					vshardrouter.InstanceInfo{},
					vshardrouter.InstanceInfo{},
				},
			},
			IsErr: true,
		},
		{
			Name: "valid",
			Source: map[vshardrouter.ReplicasetInfo][]vshardrouter.InstanceInfo{
				{Name: "rs_1", UUID: uuid.New()}: {
					vshardrouter.InstanceInfo{},
					vshardrouter.InstanceInfo{},
				},
			},
			IsErr: false,
		},
	}

	for _, tc := range tCases {
		tc := tc
		t.Run(fmt.Sprintf("is err: %v", tc.IsErr), func(t *testing.T) {
			t.Parallel()
			provider := NewProvider(tc.Source)
			if tc.IsErr {
				require.Error(t, provider.Validate())
				return
			}

			require.NoError(t, provider.Validate())
		})
	}
}
