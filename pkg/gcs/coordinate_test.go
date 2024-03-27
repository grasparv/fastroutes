package gcs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCoordinate(t *testing.T) {
	c1, err := Parse("13.388860,52.517037")
	require.NoError(t, err)
	require.Equal(t, 13.388860, c1.Latitude)
	require.Equal(t, 52.517037, c1.Longitude)

	_, err = Parse("75.32")
	require.Error(t, err)

	_, err = Parse("75.32,89a")
	require.Error(t, err)
}
