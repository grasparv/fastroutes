// Package webservice provides a small web service for retrieving
// the fastest route for several destinations
package webservice

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/grasparv/fastroutes/pkg/gcm"
	"github.com/grasparv/fastroutes/pkg/routes"
)

func Run() {
	http.HandleFunc("/routes", handleRoutes)
	http.ListenAndServe(":8080", nil)
}

// apiReply represents the message we will send back
// to the client on success
type apiReply struct {
	Source string         `json:"source"`
	Routes []routes.Route `json:"routes"`
}

func handleRoutes(w http.ResponseWriter, req *http.Request) {
	// Validate HTTP method
	if req.Method != "GET" && req.Method != "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Extract the URI query values
	u, err := url.Parse(req.RequestURI)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	values := u.Query()
	qdsts := values["dst"]
	qsrcs := values["src"]
	if len(qsrcs) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(qdsts) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse source from URI query
	src, err := gcm.Parse(qsrcs[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse destinations from URI query
	dsts := make([]gcm.Coordinate, 0, len(qdsts))
	for _, qdst := range qdsts {
		dst, err := gcm.Parse(qdst)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dsts = append(dsts, dst)
	}

	// Get the routes for source/destinations
	routes, err := routes.GetRoutes(context.Background(), src, dsts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond to client
	reply := apiReply{
		Source: src.String(),
		Routes: routes,
	}

	data, err := json.Marshal(reply)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
