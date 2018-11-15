package server

import (
	"cool-lang-features/database"
	"cool-lang-features/rpc"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type RegexHandlerMapping struct {
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

type Router struct {
	routes []RegexHandlerMapping
}

type Server struct {
	db                *database.Database
	router            *Router
	backend           string
	backendConnection net.Conn
	backendUp         bool
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
	backend := backendHostname + ":" + strconv.Itoa(backendPort)
	server := &Server{database.CreateDatabase(),
		CreateRouter(),
		backend,
		nil,
		false}
	server.connectToBackend()
	go heartbeatMonitor(server)
	return server
}

/** Attempts to connect the server to the specified backend server
 * lhs server pointer
 */
func (s *Server) connectToBackend() {
	conn, err := net.Dial("tcp", s.backend)
	s.backendConnection = conn
	s.backendUp = err == nil
}

/** Monitors heartbeat messages from the backend; if none is received, the
 * backend is assumed to be down and the monitor will attempt to reconnect
 * @param s server pointer
 */
func heartbeatMonitor(s *Server) {
	for {
		d := json.NewDecoder(s.backendConnection)
		if !s.backendUp {
			fmt.Printf("Detected failure on %s at %v\n", s.backend, time.Now())
			s.connectToBackend()
			continue
		}

		var rpcMsg rpc.RPCMapping
		err := d.Decode(&rpcMsg)
		if err != nil {
			fmt.Println("encountered error")
			fmt.Println(err)
			s.backendConnection.Close()
			s.backendUp = false
			continue
		}
		rpcType, exists := rpcMsg["type"]
		if !exists || rpcType.(string) != "Heartbeat" {
			log.Fatal("did not receive heartbeat")
		}

		fmt.Println("received heartbeat")
	}
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
