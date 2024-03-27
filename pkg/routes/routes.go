// Package routes provides functions for retrieving route information
// from the OSRM web service.
package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/grasparv/fastroutes/pkg/gcs"
)

// RouteData represents properties about a single route between two GCM
// coordinates.
type RouteData struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
}

// Route represents a route from some source to a destination.
type Route struct {
	Destination string `json:"destination"`
	RouteData
}

// apiReply represents the message response from the OSRM server
type apiReply struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Routes  []RouteData `json:"routes"`
}

// GetRoutes fetches a list of routes from the remote web service and returns
// the list sorted in order of travel duration and distance.
func GetRoutes(ctx context.Context, src gcs.Coordinate, dsts []gcs.Coordinate) (routes []Route, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	routes = make([]Route, len(dsts))

	// Serially make HTTP calls for to HTTP/1.0 server with keep-alive
	// extension.
	for i, dst := range dsts {
		var r RouteData

		keepalive := true
		if i == len(dsts)-1 {
			keepalive = false
		}

		r, err = GetRoute(ctx, src, dst, keepalive)
		if err != nil {
			return
		}

		routes[i] = Route{
			Destination: dst.String(),
			RouteData:   r,
		}
	}

	// Sort the list of routes
	slices.SortFunc(routes, func(a, b Route) int {
		if a.Duration < b.Duration {
			return -1
		} else if a.Duration > b.Duration {
			return 1
		}

		if a.Distance < b.Distance {
			return -1
		} else if a.Distance > b.Distance {
			return 1
		}

		return 0
	})

	return
}

// GetRoute fetches a single route from src to dst from the remote web service.
func GetRoute(ctx context.Context, src gcs.Coordinate, dst gcs.Coordinate, keepalive bool) (r RouteData, err error) {
	var reply apiReply
	var request *http.Request
	var response *http.Response

	// Build request
	request, err = buildRequest(ctx, src, dst, keepalive)
	if err != nil {
		return
	}

	// Send request
	response, err = http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	// Read response
	reply, err = readResponse(response.Body)
	if err != nil {
		return
	}

	// Validate API reply
	err = validateReply(reply)
	if err != nil {
		// The API response is more descriptive than the HTTP status code,
		// favouring returning this over HTTP status code.
		return
	}

	// Check HTTP status code
	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("invalid HTTP response: %s", response.Status)
		return
	}

	r = reply.Routes[0]
	return
}

// buildRequest returns a HTTP request to the OSRM web service
// for retrieving a driving route
func buildRequest(ctx context.Context, src gcs.Coordinate, dst gcs.Coordinate, keepalive bool) (*http.Request, error) {
	const method = "GET"
	const serviceUrl = "http://router.project-osrm.org/route/v1/driving"

	url := fmt.Sprintf("%s/%s;%s?overview=false", serviceUrl, src, dst)

	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	if keepalive {
		request.Header.Set("Connection", "keep-alive")
	} else {
		request.Header.Set("Connection", "close")
	}

	return request, nil
}

// readResponse returns an API reply message from the given reader
func readResponse(r io.Reader) (reply apiReply, err error) {
	var data []byte

	data, err = io.ReadAll(r)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &reply)
	if err != nil {
		return
	}

	return
}

// validateReply returns an error if the reply is invalid
func validateReply(reply apiReply) error {
	if reply.Code != "Ok" {
		return fmt.Errorf("invalid API response: %s %s", reply.Code, reply.Message)
	}

	if l := len(reply.Routes); l != 1 {
		return fmt.Errorf("invalid number of routes between source and destination: %d", l)
	}

	return nil
}
