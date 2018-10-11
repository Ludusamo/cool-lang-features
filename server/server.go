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
	db     *database.Database
	router *Router
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
func CreateServer() *Server {
	return &Server{database.CreateDatabase(), CreateRouter()}
}

/** Adds dummy data to the database for so there are a few data points
 * @lhs server pointer
 */
func (s *Server) AddDummyData() {
	s.db.AddFeature("Pattern Matching", "Pattern matching is a tool in "+
		"programming languages to process data based on its structure.")
	s.db.AddFeature("Reflection", "Reflection is a method by which a program "+
		"can achieve metaprogramming capabilities.")
}

/** Spins up the service to listen to external http requests
 * @lhs server pointer
 * @param port integer port where the service will listen from
 */
func (s *Server) Start(port int) {
	http.ListenAndServe(":"+strconv.Itoa(port), s.router)
}
