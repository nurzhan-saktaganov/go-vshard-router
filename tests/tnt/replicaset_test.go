package tnt

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/pool"
	vshardrouter "github.com/tarantool/go-vshard-router"
	"github.com/tarantool/go-vshard-router/providers/static"
)

func TestReplicsetCallAsync(t *testing.T) {
	skipOnInvalidRun(t)

	t.Parallel()

	ctx := context.Background()

	cfg := getCfg()

	router, err := vshardrouter.NewRouter(ctx, vshardrouter.Config{
		TopologyProvider: static.NewProvider(cfg),
		DiscoveryTimeout: 5 * time.Second,
		DiscoveryMode:    vshardrouter.DiscoveryModeOn,
		TotalBucketCount: totalBucketCount,
		User:             defaultTntUser,
		Password:         defaultTntPassword,
	})

	require.Nil(t, err, "NewRouter finished successfully")

	rsMap := router.RouterRouteAll()

	var rs *vshardrouter.Replicaset
	// pick random rs
	for _, v := range rsMap {
		rs = v
		break
	}

	callOpts := vshardrouter.ReplicasetCallOpts{
		PoolMode: pool.ANY,
	}

	// Tests for arglen ans response parsing
	future := rs.CallAsync(ctx, callOpts, "echo", nil)
	resp, err := future.Get()
	require.Nil(t, err, "CallAsync finished with no err on nil args")
	require.Equal(t, resp, []interface{}{}, "CallAsync returns empty arr on nil args")
	var typed interface{}
	err = future.GetTyped(&typed)
	require.Nil(t, err, "GetTyped finished with no err on nil args")
	require.Equal(t, []interface{}{}, resp, "GetTyped returns empty arr on nil args")

	const checkUpTo = 100
	for argLen := 1; argLen <= checkUpTo; argLen++ {
		args := []interface{}{}

		for i := 0; i < argLen; i++ {
			args = append(args, "arg")
		}

		future := rs.CallAsync(ctx, callOpts, "echo", args)
		resp, err := future.Get()
		require.Nilf(t, err, "CallAsync finished with no err for argLen %d", argLen)
		require.Equalf(t, args, resp, "CallAsync resp ok for argLen %d", argLen)

		var typed interface{}
		err = future.GetTyped(&typed)
		require.Nilf(t, err, "GetTyped finished with no err for argLen %d", argLen)
		require.Equal(t, args, typed, "GetTyped resp ok for argLen %d", argLen)
	}

	// Test for async execution
	timeBefore := time.Now()

	var futures = make([]*tarantool.Future, 0, len(rsMap))
	for _, rs := range rsMap {
		future := rs.CallAsync(ctx, callOpts, "sleep", []interface{}{1})
		futures = append(futures, future)
	}

	for i, future := range futures {
		_, err := future.Get()
		require.Nil(t, err, "future[%d].Get finished with no err for async test", i)
	}

	duration := time.Since(timeBefore)
	require.True(t, len(rsMap) > 1, "Async test: more than one replicaset")
	require.Less(t, duration, 1200*time.Millisecond, "Async test: requests were sent concurrently")

	// Test no timeout by default
	future = rs.CallAsync(ctx, callOpts, "sleep", []interface{}{1})
	_, err = future.Get()
	require.Nil(t, err, "CallAsync no timeout by default")

	// Test for timeout via ctx
	ctxTimeout, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	future = rs.CallAsync(ctxTimeout, callOpts, "sleep", []interface{}{1})
	_, err = future.Get()
	require.NotNil(t, err, "CallAsync timeout by context does work")

	// Test for timeout via config
	callOptsTimeout := vshardrouter.ReplicasetCallOpts{
		PoolMode: pool.ANY,
		Timeout:  500 * time.Millisecond,
	}
	future = rs.CallAsync(ctx, callOptsTimeout, "sleep", []interface{}{1})
	_, err = future.Get()
	require.NotNil(t, err, "CallAsync timeout by callOpts does work")

	future = rs.CallAsync(ctx, callOpts, "raise_luajit_error", nil)
	_, err = future.Get()
	require.NotNil(t, err, "raise_luajit_error returns error")

	future = rs.CallAsync(ctx, callOpts, "raise_client_error", nil)
	_, err = future.Get()
	require.NotNil(t, err, "raise_client_error returns error")
}

func TestReplicasetBucketsCount(t *testing.T) {
	skipOnInvalidRun(t)

	t.Parallel()

	ctx := context.Background()

	cfg := getCfg()

	router, err := vshardrouter.NewRouter(ctx, vshardrouter.Config{
		TopologyProvider: static.NewProvider(cfg),
		DiscoveryTimeout: 5 * time.Second,
		DiscoveryMode:    vshardrouter.DiscoveryModeOn,
		TotalBucketCount: totalBucketCount,
		User:             defaultTntUser,
		Password:         defaultTntPassword,
	})

	require.NoError(t, err, "NewRouter finished successfully")
	for _, rs := range router.RouterRouteAll() {
		count := uint64(0)

		count, err = rs.BucketsCount(ctx)
		require.NoError(t, err)
		require.NotEqual(t, count, 0)
	}
}
