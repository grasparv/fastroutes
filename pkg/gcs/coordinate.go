// Package gcs provides data types and functions for the Geographic coordinate system (GCS).
package gcs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

// String returns the coordinate as a comma-separated tuple.
func (c Coordinate) String() string {
	return fmt.Sprintf("%f,%f", c.Latitude, c.Longitude)
}

// Parse parses a string like 13.388860,52.517037 to a Coordinate or
// returns an error.
func Parse(s string) (c Coordinate, err error) {
	parts := strings.Split(s, ",")

	if len(parts) != 2 {
		// Note: avoid leaking the coordinate in the error message
		err = errors.New("invalid coordinate, expected comma-separated tuple")
		return
	}

	c.Latitude, err = strconv.ParseFloat(parts[0], 64)
	if err != nil {
		err = errors.New("malformed latitude, expected floating number")
		return
	}

	c.Longitude, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		err = errors.New("malformed longitude, expected floatin number")
		return
	}

	return
}
