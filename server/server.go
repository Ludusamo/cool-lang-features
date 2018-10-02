package server

import (
	"github.com/Ludusamo/cool-lang-features/database"
	"net/http"
	"strconv"
)

type Server struct {
	db     *database.Database
	router *http.ServeMux
}

func CreateServer() *Server {
	return &Server{database.CreateDatabase(), http.NewServeMux()}
}

func (s *Server) Start(port int) {
	http.ListenAndServe(":"+strconv.Itoa(port), s.router)
}
