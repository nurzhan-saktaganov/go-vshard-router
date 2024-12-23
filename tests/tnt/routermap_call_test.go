package tnt

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	vshardrouter "github.com/tarantool/go-vshard-router"
	"github.com/tarantool/go-vshard-router/providers/static"
)

func TestRouterMapCall(t *testing.T) {
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

	callOpts := vshardrouter.CallOpts{}

	const arg = "arg1"

	// Enusre that RouterMapCallRWImpl works at all
	echoArgs := []interface{}{arg}
	resp, err := router.RouterMapCallRWImpl(ctx, "echo", echoArgs, callOpts)
	require.NoError(t, err, "RouterMapCallRWImpl echo finished with no err")

	for k, v := range resp {
		require.Equalf(t, arg, v, "RouterMapCallRWImpl value ok for %v", k)
	}

	echoArgs = []interface{}{1}
	respInt, err := vshardrouter.RouterMapCallRW[int](router, ctx, "echo", echoArgs, vshardrouter.RouterMapCallRWOptions{})
	require.NoError(t, err, "RouterMapCallRW[int] echo finished with no err")
	for k, v := range respInt {
		require.Equalf(t, 1, v, "RouterMapCallRW[int] value ok for %v", k)
	}

	// RouterMapCallRWImpl returns only one value
	echoArgs = []interface{}{arg, "arg2"}
	resp, err = router.RouterMapCallRWImpl(ctx, "echo", echoArgs, callOpts)
	require.NoError(t, err, "RouterMapCallRWImpl echo finished with no err")

	for k, v := range resp {
		require.Equalf(t, arg, v, "RouterMapCallRWImpl value ok for %v", k)
	}

	// RouterMapCallRWImpl returns nil when no return value
	noArgs := []interface{}{}
	resp, err = router.RouterMapCallRWImpl(ctx, "echo", noArgs, callOpts)
	require.NoError(t, err, "RouterMapCallRWImpl echo finished with no err")

	for k, v := range resp {
		require.Equalf(t, nil, v, "RouterMapCallRWImpl value ok for %v", k)
	}

	// Ensure that RouterMapCallRWImpl sends requests concurrently
	const sleepToSec int = 1
	sleepArgs := []interface{}{sleepToSec}

	start := time.Now()
	_, err = router.RouterMapCallRWImpl(ctx, "sleep", sleepArgs, vshardrouter.CallOpts{
		Timeout: 2 * time.Second, // because default timeout is 0.5 sec
	})
	duration := time.Since(start)

	require.NoError(t, err, "RouterMapCallRWImpl sleep finished with no err")
	require.Greater(t, len(cfg), 1, "There are more than one replicasets")
	require.Less(t, duration, 1200*time.Millisecond, "Requests were send concurrently")

	// RouterMapCallRWImpl returns err on raise_luajit_error
	_, err = router.RouterMapCallRWImpl(ctx, "raise_luajit_error", noArgs, callOpts)
	require.NotNil(t, err, "RouterMapCallRWImpl raise_luajit_error finished with error")

	// RouterMapCallRWImpl invalid usage
	_, err = router.RouterMapCallRWImpl(ctx, "echo", nil, callOpts)
	require.NotNil(t, err, "RouterMapCallRWImpl with nil args finished with error")

	// Ensure that RouterMapCallRWImpl doesn't work when it mean't to
	for k := range cfg {
		errs := router.RemoveReplicaset(ctx, k.UUID)
		require.Emptyf(t, errs, "%s successfully removed from router", k.UUID)
		break
	}

	_, err = router.RouterMapCallRWImpl(ctx, "echo", echoArgs, callOpts)
	require.NotNilf(t, err, "RouterMapCallRWImpl failed on not full cluster")
}
