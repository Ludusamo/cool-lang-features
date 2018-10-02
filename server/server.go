package server

import (
	"github.com/Ludusamo/cool-lang-features/database"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Server struct {
	db     *database.Database
	router *mux.Router
}

func CreateServer() *Server {
	return &Server{database.CreateDatabase(), mux.NewRouter()}
}

func (s *Server) Start(port int) {
	http.ListenAndServe(":"+strconv.Itoa(port), s.router)
}
