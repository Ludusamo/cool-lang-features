package routes

import (
	"github.com/Ludusamo/cool-lang-features/server"
	"io"
	"net/http"
)

func (s *Server) RegisterHandlers() {
	s.router.HandleFunc("/", s.homeHandler())
	s.router.HandleFunc("/api", s.apiHandler())
}

func (s *Server) homeHandler() {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, "Hello Web!")
	}
}

func (s *Server) apiHandler() {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, "api")
	}
}
