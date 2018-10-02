package server

import (
	"io"
	"net/http"
    "encoding/json"
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
		}
	}
}
