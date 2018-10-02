package server

import (
	"github.com/Ludusamo/cool-lang-features/database"
	"net/http"
)

type Server struct {
	db     *database.Database
	router *http.ServeMux
}

func CreateServer() *Server {
	return &Server{database.CreateDatabase(), http.NewServeMux()}
}
