package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

/** Registers handlers for all the different routes on the server
 * @lhs server pointer
 */
func (s *Server) RegisterHandlers() {
	s.router.HandleFunc("/api", s.apiHandler())
	s.router.HandleFunc("/api/feature", s.featuresHandler())
	s.router.HandleFunc("/api/feature/{id:[0-9]+}", s.featureHandler())
	// Needs to be placed last. These are evaluated in the order they are added
	s.router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/index.html")
	})
}

/** Handler for the /api route
 * @lhs server pointer
 * @return Handler function for router
 */
func (s *Server) apiHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, "api")
	}
}

/** Handler for /api/feature for adding and retrieving features
 * @lhs server pointer
 * @return Handler function for router
 */
func (s *Server) featuresHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			features := s.db.GetFeatures()
			json.NewEncoder(w).Encode(features)
		case "POST":
			var featurePost struct {
				Name        string
				Description string
			}
			json.NewDecoder(r.Body).Decode(&featurePost)
			feat, err := s.db.AddFeature(
				featurePost.Name,
				featurePost.Description)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(err)
			} else {
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(feat)
			}
		}
	}
}

/** Handler for /api/feature/{id} for modifying, deleting and retrieving
 * features based on id
 * @lhs server pointer
 * @return Handler function for router
 */
func (s *Server) featureHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, parseErr := strconv.Atoi(vars["id"])
		if parseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(parseErr)
			return
		}
		switch r.Method {
		case "GET":
			feat, err := s.db.GetFeature(id)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(err)
			} else {
				json.NewEncoder(w).Encode(feat)
			}
		case "PATCH":
			var featurePatch struct {
				Name        string
				Description string
			}
			json.NewDecoder(r.Body).Decode(&featurePatch)
			feat, err := s.db.ModifyFeature(id,
				featurePatch.Name,
				featurePatch.Description)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(err)
			} else {
				json.NewEncoder(w).Encode(feat)
			}
		case "DELETE":
			s.db.DeleteFeature(id)
			w.WriteHeader(http.StatusOK)
		}
	}
}
