package routes

import (
	"context"
	"testing"

	"github.com/grasparv/fastroutes/pkg/gcs"
	"github.com/stretchr/testify/require"
)

func TestRoute(t *testing.T) {
	ctx := context.Background()
	src := gcs.Coordinate{13.388860, 52.517037}
	dst := gcs.Coordinate{13.397634, 52.529407}
	r, err := GetRoute(ctx, src, dst, false)
	require.NoError(t, err)
	require.Equal(t, float64(1886.8), r.Distance)
	require.Equal(t, float64(260.3), r.Duration)
}

func TestRoutes(t *testing.T) {
	ctx := context.Background()
	src := gcs.Coordinate{13.388860, 52.517037}
	dstNear := gcs.Coordinate{13.397634, 52.529407}
	dstFar := gcs.Coordinate{14.397634, 53.529407}
	r, err := GetRoutes(ctx, src, []gcs.Coordinate{dstFar, dstNear})
	require.NoError(t, err)
	require.Equal(t, 2, len(r))
	require.Less(t, r[0].Distance, r[1].Distance)
	require.Less(t, r[0].Duration, r[1].Duration)
	require.Equal(t, dstNear.String(), r[0].Destination)
	require.Equal(t, dstFar.String(), r[1].Destination)
}
