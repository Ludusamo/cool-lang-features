package server

import (
	"io"
	"net/http"
    "encoding/json"
	"github.com/gorilla/mux"
	"strconv"
)

func (s *Server) RegisterHandlers() {
	s.router.HandleFunc("/", s.homeHandler())
	s.router.HandleFunc("/api", s.apiHandler())
	s.router.HandleFunc("/api/feature", s.featuresHandler())
}

func (s *Server) homeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, "Hello Web!")
	}
}

func (s *Server) apiHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, "api")
	}
}

func (s *Server) featuresHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			features := s.db.GetFeatures()
            json.NewEncoder(w).Encode(features)
        case "POST":
            var featurePost struct {
                Name string
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
                json.NewEncoder(w).Encode(feat)
            }
        }
	}
}

func (s *Server) featureHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
		switch r.Method {
		case "GET":
            id, parseErr := strconv.Atoi(vars["id"])
            if parseErr != nil {
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(parseErr)
                return
            }
            feat, err := s.db.GetFeature(id)
            if err != nil {
                w.WriteHeader(http.StatusNotFound)
                json.NewEncoder(w).Encode(err)
            } else {
                json.NewEncoder(w).Encode(feat)
            }
        case "PATCH":
        case "DELETE":
		}
	}
}
