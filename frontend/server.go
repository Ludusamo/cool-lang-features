package server

import (
	"cool-lang-features/database"
	"net/http"
	"regexp"
	"strconv"
)

type RegexHandlerMapping struct {
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

type Router struct {
	routes []RegexHandlerMapping
}

type Server struct {
	db      *database.Database
	router  *Router
	backend string
}

/** Creates an empty router object
 * @return pointer to created router
 */
func CreateRouter() *Router {
	return &Router{make([]RegexHandlerMapping, 0)}
}

/** Creates an empty server object
 * @return pointer to created server
 */
func CreateServer(backendHostname string, backendPort int) *Server {
	return &Server{database.CreateDatabase(),
		CreateRouter(),
		backendHostname + ":" + strconv.Itoa(backendPort)}
}

/** Spins up the service to listen to external http requests
 * @lhs server pointer
 * @param port integer port where the service will listen from
 * @param backendHostname hostname of the backend
 * @param backendPort port of the backend
 */
func (s *Server) Start(port int) {
	http.ListenAndServe(":"+strconv.Itoa(port), s.router)
}
